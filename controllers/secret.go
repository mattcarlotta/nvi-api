package controllers

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
	"gorm.io/gorm"
)

func GetSecretBySecretID(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	id := c.Params("id")
	if err := utils.Validate().Var(id, "required,uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetSecretInvalidID))
	}

	var secret models.Secret
	if err := db.Preload("Environments").First(
		&secret, "id=? AND user_id=?", utils.MustParseUUID(id), userSessionID,
	).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.GetSecretNonExistentID))
	}

	return c.Status(fiber.StatusOK).JSON(secret)
}

func GetSecretsByEnvironmentID(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	id := c.Params("id")
	if err := utils.Validate().Var(id, "required,uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetSecretsByEnvInvalidID))
	}

	parsedEnvID := utils.MustParseUUID(id)
	var secrets []models.SecretResult
	if err := db.Raw(
		utils.FindSecretsByEnvIDQuery, userSessionID, utils.GenerateJSONIDString(parsedEnvID),
	).Scan(&secrets).Error; err != nil || len(secrets) == 0 {
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
		}
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.GetSecretsByEnvNonExistentID))
	}

	return c.Status(fiber.StatusOK).JSON(secrets)
}

func CreateSecret(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	data := new(models.ReqCreateSecret)
	if err := c.BodyParser(data); err != nil {
		log.Printf("Failing at JSON parsing: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.CreateSecretInvalidBody))
	}

	if err := utils.Validate().Struct(data); err != nil {
		log.Printf("Failing at JSON validation: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.CreateSecretInvalidBody))
	}

	environmentIDs, err := utils.ParseUUIDs(data.EnvironmentIDs)
	if err != nil {
		log.Printf("Failing at UUID parsing: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	var environments []models.Environment
	if err := db.Find(
		&environments, "id IN ? AND user_id=?", environmentIDs, userSessionID,
	).Error; err != nil || len(environments) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.CreateSecretNonExistentEnv))
	}

	// var secrets []models.SecretResult
	// if err := db.Raw(
	//      utils.GenerateFindSecretByEnvIDsQuery, userSessionID, data.Key,
	// ).Scan(&secrets).Error; err != nil || len(secrets) > 0 {
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
	if err := db.Preload("Environments", "ID in ?", environmentIDs).Find(
		&secrets, "key=? AND user_id=?", data.Key, userSessionID,
	).Error; err == nil {
		duplicates := models.GetDupKeyinEnvs(&secrets)
		if len(duplicates) > 0 {
			return c.Status(fiber.StatusConflict).JSON(utils.JSONError(utils.CreateSecretKeyAlreadyExists))
		}
	}

	newSecret := models.Secret{Key: data.Key, Value: []byte(data.Value), UserID: userSessionID, Environments: environments}
	if err := db.Create(&newSecret).Error; err != nil {
		log.Printf("Failing at creating Secret: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	return c.Status(fiber.StatusCreated).SendString(fmt.Sprintf("Successfully created %s!", data.Key))
}

func DeleteSecret(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	id := c.Params("id")
	if err := utils.Validate().Var(id, "required,uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.DeleteSecretInvalidID))
	}

	var secret models.Secret
	if err := db.Where(
		&models.Environment{ID: utils.MustParseUUID(id), UserID: userSessionID},
	).First(&secret).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.DeleteSecretNonExistentID))
	}

	if err := db.Delete(&secret).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	return c.Status(fiber.StatusCreated).SendString(fmt.Sprintf("Successfully deleted the %s secret!", secret.Key))
}

func UpdateSecret(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	data := new(models.ReqUpdateSecret)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide valid environments and a secret key name with content!"},
		)
	}

	if err := utils.Validate().Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "You must provide a valid secret id, one or more environment ids and a secret key name with content!",
			},
		)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		parsedID := utils.MustParseUUID(data.ID)

		var secret models.Secret
		if err := tx.Where(&models.Secret{ID: parsedID, UserID: userSessionID}).First(&secret).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{"error": "The provided secret doesn't appear to exist!"},
			)
		}

		environmentIDs, err := utils.ParseUUIDs(data.EnvironmentIDs)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
		}

		var environments []models.Environment
		if err := tx.Find(
			&environments, "id IN ? AND user_id=?", environmentIDs, userSessionID,
		).Error; err != nil || len(environments) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{"error": "The provided environment doesn't appear to exist!"},
			)
		}

		var secrets []models.Secret
		if err := tx.Preload(
			"Environments", "ID in ?", environmentIDs,
		).Not(
			"id", parsedID,
		).Find(
			&secrets, "key=? AND user_id=?", data.Key, userSessionID,
		).Error; err != nil || len(secrets) > 0 {
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
			}

			duplicates := models.GetDupKeyinEnvs(&secrets)
			if len(duplicates) > 0 {
				return c.Status(fiber.StatusConflict).JSON(
					fiber.Map{"error": fmt.Sprintf(
						"The key '%s' already exists for the following selected environments: %s.", data.Key, duplicates),
					},
				)
			}

		}

		newValue, err := utils.CreateEncryptedText([]byte(data.Key))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
		}

		if err = tx.Model(&secret).Association("Environments").Replace(environments); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
		}

		if err = tx.Model(&secret).Updates(&models.Secret{Key: data.Key, Value: newValue}).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
		}

		return c.Status(fiber.StatusCreated).SendString(fmt.Sprintf("Successfully updated the %s secret!", data.Key))
	})
}
