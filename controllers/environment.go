package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
	"gorm.io/gorm"
)

func GetAllEnvironments(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	var environments []models.Environment
	db.Where(&models.Environment{UserID: userSessionID}).Find(&environments)

	return c.Status(fiber.StatusOK).JSON(environments)
}

func GetEnvironmentByID(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	id := c.Params("id")
	if err := utils.Validate().Var(id, "required,uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetEnvironmentInvalidID))
	}

	var environment models.Environment
	if err := db.Where(
		&models.Environment{ID: utils.MustParseUUID(id), UserID: userSessionID},
	).First(&environment).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.GetEnvironmentNonExistentID))
	}

	return c.Status(fiber.StatusOK).JSON(environment)
}

func CreateEnvironment(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	envName := c.Params("name")
	if err := utils.Validate().Var(envName, "required,envname,lte=255"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.CreateEnvironmentInvalidName))
	}

	newEnv := models.Environment{Name: envName, UserID: userSessionID}
	var environment models.Environment
	if err := db.Where(&newEnv).First(&environment).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(utils.JSONError(utils.CreateEnvironmentNameTaken))
	}

	// TODO(carlotta): add a limit to how many environments can be created per account

	if err := db.Create(&newEnv).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	return c.Status(fiber.StatusCreated).SendString(fmt.Sprintf("Successfully created a(n) %s environment!", envName))
}

func DeleteEnvironment(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	id := c.Params("id")
	if err := utils.Validate().Var(id, "required,uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.DeleteEnvironmentInvalidID))
	}

	return db.Transaction(func(tx *gorm.DB) error {
		parsedID := utils.MustParseUUID(id)

		var environment models.Environment
		if err := tx.Where(
			&models.Environment{ID: parsedID, UserID: userSessionID},
		).First(&environment).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.DeleteEnvironmentNonExistentID))
		}

		var secrets []models.Secret
		if err := tx.Preload("Environments").Not("id", parsedID).Find(
			&secrets, "user_id=?", userSessionID,
		).Error; err != nil || len(secrets) > 0 {
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
			}

			for _, secret := range secrets {
				for _, env := range secret.Environments {
					if env.ID == environment.ID {
						tx.Delete(&secret)
					}
				}
			}
		}

		if err := tx.Delete(&environment).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
		}

		return c.Status(fiber.StatusOK).SendString(
			fmt.Sprintf("Successfully deleted the %s environment!", environment.Name),
		)
	})

}

func UpdateEnvironment(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	var data models.ReqUpdateEnv
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.UpdateEnvironmentInvalidBody))
	}

	if err := utils.Validate().Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.UpdateEnvironmentInvalidBody))
	}

	var environment models.Environment
	if err := db.Where(
		&models.Environment{ID: utils.MustParseUUID(data.ID), UserID: userSessionID},
	).First(&environment).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.UpdateEnvironmentNonExistentID))
	}

	if err := db.Model(&environment).Update("name", &data.UpdatedName).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	return c.Status(fiber.StatusOK).SendString(
		fmt.Sprintf("Successfully updated the environment name to %s!", data.UpdatedName),
	)
}
