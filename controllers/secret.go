package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
	"gorm.io/gorm"
)

type ReqSecret struct {
	ID             string   `json:"id"`
	EnvironmentIds []string `json:"environmentIds"`
	Name           string   `json:"name"`
	Content        string   `json:"content"`
}

// func findDupEnvNames(secrets *[]models.Secret, ids []uuid.UUID) string {
// 		var envNames string
// 		for _, secret := range *secrets {
// 			for _, env := range secret.Environments {
// 				for _, id := range ids {
// 					if env.ID == id {
// 						if len(envNames) == 0 {
// 							envNames += env.Name
// 						} else {
// 							envNames += fmt.Sprintf(", %s", env.Name)
// 						}
// 					}
// 				}
// 			}
// 		}

//         return envNames
// }

func getDupKeyinEnvs(secrets *[]models.Secret) string {
	var envNames string
	for _, secret := range *secrets {
		for _, env := range secret.Environments {
			if len(envNames) == 0 {
				envNames += env.Name
			} else {
				envNames += fmt.Sprintf(", %s", env.Name)
			}
		}
	}

	return envNames
}

func CreateSecret(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := utils.GetSessionId(c)

	data := new(ReqSecret)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide valid environments and a secret key name with content!"},
		)
	}

	// TODO(carlotta): Add field validations for "environmentIds," "name," and "content"
	if len(data.EnvironmentIds) == 0 || len(data.Name) == 0 || len(data.Content) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide valid environments and a secret key name with content!"},
		)
	}

	environmentIds, err := utils.ParseUUIDs(data.EnvironmentIds)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": err.Error()},
		)
	}

	var environments []models.Environment
	if err := db.Find(
		&environments, "id IN ? AND user_id=?", environmentIds, userSessionId,
	).Error; err != nil || len(environments) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "The provided environment doesn't appear to exist!"},
		)
	}

	var secrets []models.Secret
	if err := db.Preload("Environments", "ID in ?", environmentIds).Find(
		&secrets, "name=? AND user_id=?", data.Name, userSessionId,
	).Error; err == nil {
		envNames := getDupKeyinEnvs(&secrets)
		if len(envNames) > 0 {
			return c.Status(fiber.StatusConflict).JSON(
				fiber.Map{"error": fmt.Sprintf(
					"The key '%s' already exists for the following selected environments: %s", data.Name, envNames),
				},
			)
		}
	}

	newSecret := models.Secret{
		Name:         data.Name,
		Content:      []byte(data.Content),
		UserId:       userSessionId,
		Environments: environments,
	}
	if err := db.Create(&newSecret).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).SendString(
		fmt.Sprintf("Successfully created %s!", data.Name),
	)
}

func DeleteSecret(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := utils.GetSessionId(c)

	parsedId, err := utils.ParseUUID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid secret id!"},
		)
	}

	var secret models.Secret
	if err := db.Where(
		&models.Environment{ID: parsedId, UserId: userSessionId},
	).First(&secret).Error; err != nil {
		return c.Status(fiber.StatusOK).JSON(
			fiber.Map{"error": "The provided secret doesn't appear to exist!"},
		)
	}

	if err := db.Delete(&secret).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).SendString(
		fmt.Sprintf("Successfully deleted the %s secret!", secret.Name),
	)
}

func UpdateSecret(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := utils.GetSessionId(c)

	data := new(ReqSecret)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide valid environments and a secret key name with content!"},
		)
	}

	// TODO(carlotta): Add field validations for "id", "environmentIds," "name," and "content"
	if len(data.ID) == 0 || len(data.EnvironmentIds) == 0 || len(data.Name) == 0 || len(data.Content) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "You must provide a valid secret id, one or more environment ids and a secret key name with content!",
			},
		)
	}

	parsedId, err := utils.ParseUUID(data.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": err.Error()},
		)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		var secret models.Secret
		if err := tx.Where(&models.Secret{ID: parsedId, UserId: userSessionId}).First(&secret).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{"error": "The provided secret doesn't appear to exist!"},
			)
		}

		environmentIds, err := utils.ParseUUIDs(data.EnvironmentIds)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{"error": err.Error()},
			)
		}

		var environments []models.Environment
		if err := tx.Find(
			&environments, "id IN ? AND user_id=?", environmentIds, userSessionId,
		).Error; err != nil || len(environments) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{"error": "The provided environment doesn't appear to exist!"},
			)
		}

		var secrets []models.Secret
		if err := tx.Preload(
			"Environments", "ID in ?", environmentIds,
		).Not("id", parsedId).Find(
			&secrets, "name=? AND user_id=?", data.Name, userSessionId,
		).Error; err != nil || len(secrets) != 0 {
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			envNames := getDupKeyinEnvs(&secrets)
			if len(envNames) > 0 {
				return c.Status(fiber.StatusConflict).JSON(
					fiber.Map{"error": fmt.Sprintf(
						"The key '%s' already exists for the following selected environments: %s.", data.Name, envNames),
					},
				)
			}

		}

		newContent, err := utils.CreateEncryptedText([]byte(data.Content))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err = tx.Model(&secret).Association("Environments").Replace(environments); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err = tx.Model(&secret).Updates(&models.Secret{Name: data.Name, Content: newContent}).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).SendString(
			fmt.Sprintf("Successfully updated the %s secret!", data.Name),
		)
	})
}
