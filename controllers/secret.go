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
	Key            string   `json:"key"`
	Value          string   `json:"value"`
}

type EnvironmentName struct {
	Name string `json:"name"`
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

func GetSecretBySecretId(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := utils.GetSessionId(c)

	parsedSecretId, err := utils.ParseUUID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id!"},
		)
	}

	var secret models.Secret
	if err := db.Preload("Environments").First(
		&secret, "id=? AND user_id=?", parsedSecretId, userSessionId,
	).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{"error": "Unable to locate a secret with that id!"},
		)
	}

	return c.Status(fiber.StatusOK).JSON(secret)
}

func GetSecretsByEnvironmentId(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := utils.GetSessionId(c)

	parsedEnvId, err := utils.ParseUUID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id!"},
		)
	}

	var secrets []models.SecretResult
	if err := db.Raw(`
        SELECT * 
        FROM (
	        SELECT 
		        s.id,
		        s.user_id,
		        s.key,
		        s.value,
		        s.created_at,
		        s.updated_at,
		        jsonb_agg(envs) as environments
	        FROM secrets s
	        JOIN environment_secrets es ON s.id = es.secret_id
	        JOIN environments envs on es.environment_id = envs.id
	        WHERE s.user_id = ?
	        GROUP BY s.id
        ) r
        WHERE r.environments @> ?;
    `, userSessionId, `[{"id":"`+parsedEnvId.String()+`"}]`).Scan(&secrets).Error; err != nil {
		fmt.Printf("Failed to load secrets with %s: %s", parsedEnvId, err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Failed to locate any secrets with that id."},
		)
	}

	return c.Status(fiber.StatusOK).JSON(secrets)
}

func CreateSecret(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := utils.GetSessionId(c)

	data := new(ReqSecret)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide valid environments and a secret key with a value!"},
		)
	}

	// TODO(carlotta): Add field validations for "environmentIds," "key," and "value"
	if len(data.EnvironmentIds) == 0 || len(data.Key) == 0 || len(data.Value) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide valid environments and a secret key name with content!"},
		)
	}

	environmentIds, err := utils.ParseUUIDs(data.EnvironmentIds)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var environments []models.Environment
	if err := db.Find(
		&environments, "id IN ? AND user_id=?", environmentIds, userSessionId,
	).Error; err != nil || len(environments) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{"error": "The provided environments don't appear to exist!"},
		)
	}

	// var queryEnvironments string
	// for _, value := range environmentIds {
	// 	if len(queryEnvironments) == 0 {
	// 		queryEnvironments += `r.environments @> '[{"id": "` + value.String() + `"}]'`
	// 	} else {
	// 		queryEnvironments += `OR r.environments @> '[{"id": "` + value.String() + `"}]'`
	// 	}
	// }

	// RAWSQL := `
	//        SELECT *
	//        FROM (
	//         SELECT
	// 	        s.id,
	// 	        s.user_id,
	// 	        s.key,
	// 	        s.value,
	// 	        s.created_at,
	// 	        s.updated_at,
	// 	        jsonb_agg(envs) as environments
	//         FROM secrets s
	//         JOIN environment_secrets es ON s.id = es.secret_id
	//         JOIN environments envs on es.environment_id = envs.id
	//         WHERE s.user_id = ?
	//         GROUP BY s.id
	//        ) r
	//        WHERE `

	// RAWSQL += "(" + queryEnvironments + ") AND r.key = ?"

	// var secrets []models.SecretResult
	// if err := db.Raw(RAWSQL, userSessionId, data.Key).Scan(&secrets).Error; err != nil || len(secrets) > 0 {
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	// 	}
	//
	// 	return c.Status(fiber.StatusConflict).JSON(
	// 		fiber.Map{"error": fmt.Sprintf(
	// 			"The key '%s' already exists in one or more of the selected environments!", data.Key),
	// 		},
	// 	)
	// 	// }
	// }

	var secrets []models.Secret
	if err := db.Preload("Environments", "ID in ?", environmentIds).Find(
		&secrets, "key=? AND user_id=?", data.Key, userSessionId,
	).Error; err == nil {
		envNames := getDupKeyinEnvs(&secrets)
		if len(envNames) > 0 {
			return c.Status(fiber.StatusConflict).JSON(
				fiber.Map{"error": fmt.Sprintf(
					"The key '%s' already exists for the following selected environments: %s", data.Key, envNames),
				},
			)
		}
	}

	newSecret := models.Secret{Key: data.Key, Value: []byte(data.Value), UserId: userSessionId, Environments: environments}
	if err := db.Create(&newSecret).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).SendString(fmt.Sprintf("Successfully created %s!", data.Key))
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

	return c.Status(fiber.StatusCreated).SendString(fmt.Sprintf("Successfully deleted the %s secret!", secret.Key))
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
	if len(data.ID) == 0 || len(data.EnvironmentIds) == 0 || len(data.Key) == 0 || len(data.Key) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "You must provide a valid secret id, one or more environment ids and a secret key name with content!",
			},
		)
	}

	parsedId, err := utils.ParseUUID(data.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
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
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
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
		).Not(
			"id", parsedId,
		).Find(
			&secrets, "key=? AND user_id=?", data.Key, userSessionId,
		).Error; err != nil || len(secrets) != 0 {
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			envNames := getDupKeyinEnvs(&secrets)
			if len(envNames) > 0 {
				return c.Status(fiber.StatusConflict).JSON(
					fiber.Map{"error": fmt.Sprintf(
						"The key '%s' already exists for the following selected environments: %s.", data.Key, envNames),
					},
				)
			}

		}

		newValue, err := utils.CreateEncryptedText([]byte(data.Key))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err = tx.Model(&secret).Association("Environments").Replace(environments); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err = tx.Model(&secret).Updates(&models.Secret{Key: data.Key, Value: newValue}).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).SendString(fmt.Sprintf("Successfully updated the %s secret!", data.Key))
	})
}
