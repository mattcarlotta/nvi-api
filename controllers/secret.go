package controllers

import (
	"fmt"

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

func SearchForSecretsByEnvironmentIDAndSecretKey(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	key := c.Query("key")
	if err := utils.Validate().Var(key, "required,gte=2,lte=255"); err != nil {
		return c.Status(fiber.StatusBadRequest).Send(nil)
		// return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetProjectInvalidName))
	}
	environmentID := c.Query("environmentID")
	if err := utils.Validate().Var(environmentID, "required,uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetEnvironmentInvalidProjectID))
	}

	parsedEnvID := utils.MustParseUUID(environmentID)

	var secrets []models.Secret
	if err := db.Raw(
		utils.FindSecretsByEnvIDAndSecretKeyQuery, userSessionID, "%"+key+"%", utils.GenerateJSONIDString(parsedEnvID),
	).Scan(&secrets).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	return c.Status(fiber.StatusOK).JSON(secrets)
}

func CreateSecret(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	var data models.ReqCreateSecret
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.CreateSecretInvalidBody))
	}

	if err := utils.Validate().Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.CreateSecretInvalidBody))
	}

	projectID := utils.MustParseUUID(data.ProjectID)

	var project models.Project
	if err := db.Where(
		"id=? AND user_id=?", projectID, userSessionID,
	).First(&project).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.CreateSecretNonExistentProject))
	}

	environmentIDs, err := utils.ParseUUIDs(data.EnvironmentIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	var environments []models.Environment
	if err := db.Find(
		&environments, "id IN ? AND project_id=? AND user_id=?", environmentIDs, projectID, userSessionID,
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
	if err := db.Preload("Environments", "id in ? AND project_id=?", environmentIDs, projectID).Find(
		&secrets, "key=? AND user_id=?", data.Key, userSessionID,
	).Error; err == nil {
		duplicates := models.GetDupKeyinEnvs(&secrets)
		if len(duplicates) > 0 {
			return c.Status(fiber.StatusConflict).JSON(utils.JSONError(utils.CreateSecretKeyAlreadyExists))
		}
	}

	newSecret := models.Secret{Key: data.Key, Value: []byte(data.Value), UserID: userSessionID, Environments: environments}
	if err := db.Create(&newSecret).Error; err != nil {
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

	var data models.ReqUpdateSecret
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.UpdateSecretInvalidBody))
	}

	if err := utils.Validate().Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.UpdateSecretInvalidBody))
	}

	return db.Transaction(func(tx *gorm.DB) error {
		parsedID := utils.MustParseUUID(data.ID)

		var secret models.Secret
		if err := tx.Where(&models.Secret{ID: parsedID, UserID: userSessionID}).First(&secret).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.UpdateSecretInvalidID))
		}

		environmentIDs, err := utils.ParseUUIDs(data.EnvironmentIDs)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
		}

		var environments []models.Environment
		if err := tx.Find(
			&environments, "id IN ? AND user_id=?", environmentIDs, userSessionID,
		).Error; err != nil || len(environments) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.UpdateSecretNonExistentEnv))
		}

		var secrets []models.Secret
		if err := tx.Preload(
			"Environments", "ID in ?", environmentIDs,
		).Not(
			"id", parsedID,
		).Find(
			&secrets, "key=? AND user_id=?", data.Key, userSessionID,
		).Error; err != nil || len(secrets) != 0 {
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
			}

			duplicates := models.GetDupKeyinEnvs(&secrets)
			if len(duplicates) > 0 {
				return c.Status(fiber.StatusConflict).JSON(utils.JSONError(utils.UpdateSecretKeyAlreadyExists))
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

		return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Successfully updated the %s secret!", data.Key))
	})
}
