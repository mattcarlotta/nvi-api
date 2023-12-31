package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
	"gorm.io/gorm"
)

func GetSecretsByProjectAndEnvironmentName(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	projectName := c.Query("project")
	if err := utils.Validate().Var(projectName, "required,name,lte=255"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetProjectInvalidName))
	}

	environmentName := c.Query("environment")
	if err := utils.Validate().Var(environmentName, "required,name,lte=255"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetEnvironmentInvalidName))
	}

	var project models.Project
	if err := db.Where(
		&models.Project{Name: projectName, UserID: userSessionID},
	).First(&project).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.GetProjectNonExistentName))
	}

	var environment models.Environment
	if err := db.Where(
		&models.Environment{Name: environmentName, ProjectID: project.ID, UserID: userSessionID},
	).First(&environment).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.GetEnvironmentNonExistentName))
	}

	var secrets []models.SecretResult
	if err := db.Raw(
		utils.FindSecretsByEnvIDQuery, userSessionID, utils.GenerateJSONIDString(environment.ID),
	).Scan(&secrets).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	var environments []models.Environment
	db.Where(&models.Environment{UserID: userSessionID, ProjectID: project.ID}).Find(&environments)

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"environment":  environment,
			"environments": environments,
			"project":      project,
			"secrets":      secrets,
		},
	)
}

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

	decryptedValue, err := utils.DecryptSecretValue(secret.Value, secret.Nonce)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	var environmentIDs []uuid.UUID
	for _, s := range secret.Environments {
		environmentIDs = append(environmentIDs, s.ID)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{"environmentIDs": environmentIDs, "key": secret.Key, "value": string(decryptedValue)},
	)
}

func GetSecretsByEnvironmentID(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	id := c.Params("id")
	if err := utils.Validate().Var(id, "required,uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetSecretsByEnvInvalidID))
	}

	parsedEnvID := utils.MustParseUUID(id)

	var environment models.Environment
	if err := db.Where(&models.Environment{ID: parsedEnvID}).First(&environment).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.GetSecretsByEnvNonExistentID))
	}

	var secrets []models.SecretResult
	if err := db.Raw(
		utils.FindSecretsByEnvIDQuery, userSessionID, utils.GenerateJSONIDString(environment.ID),
	).Scan(&secrets).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	return c.Status(fiber.StatusOK).JSON(secrets)
}

func SearchForSecretsByEnvironmentIDAndSecretKey(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	key := c.Query("key")
	if err := utils.Validate().Var(key, "required,gte=2,lte=255"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.SearchForSecretsByEnvAndSecretInvalidKey))
	}

	environmentID := c.Query("environmentID")
	if err := utils.Validate().Var(environmentID, "required,uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetEnvironmentInvalidID))
	}

	parsedEnvID := utils.MustParseUUID(environmentID)

	var secrets []models.SecretResult
	if err := db.Raw(
		utils.FindSecretsByEnvIDAndSecretKeyQuery,
		userSessionID,
		"%"+key+"%",
		utils.GenerateJSONIDString(parsedEnvID),
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
		&models.Project{ID: projectID, UserID: userSessionID},
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
	if err := db.Preload(
		"Environments", "id in ? AND project_id=?", environmentIDs, projectID,
	).Find(
		&secrets, "key=? AND user_id=?", data.Key, userSessionID,
	).Error; err == nil {
		duplicates := models.GetDupKeyinEnvs(&secrets)
		if len(duplicates) > 0 {
			return c.Status(fiber.StatusConflict).JSON(utils.JSONError(utils.CreateSecretKeyAlreadyExists))
		}
	}

	// TODO(carlotta): add secret limitations per environment?
	// This may be costly to perform if multiple environments are selected.

	newSecret := models.Secret{
		Key:          data.Key,
		Value:        []byte(data.Value),
		UserID:       userSessionID,
		Environments: environments,
	}
	if err := db.Create(&newSecret).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	return c.Status(fiber.StatusCreated).JSON(newSecret)
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

	return c.Status(fiber.StatusCreated).SendString(fmt.Sprintf("Successfully removed the %s secret!", secret.Key))
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

		newValue, newNonce, err := utils.CreateEncryptedSecretValue([]byte(data.Value))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
		}

		if err = tx.Model(&secret).Association("Environments").Replace(environments); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
		}

		updatedSecret := models.Secret{
			Key:   data.Key,
			Nonce: newNonce,
			Value: newValue,
		}
		if err = tx.Model(&secret).Updates(&updatedSecret).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
		}

		return c.Status(fiber.StatusOK).JSON(secret)
	})
}
