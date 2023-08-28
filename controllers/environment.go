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
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id!"},
		)
	}

	var environment models.Environment
	if err := db.Where(
		&models.Environment{ID: utils.MustParseUUID(id), UserID: userSessionID},
	).First(&environment).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "The provided environment doesn't appear to exist!"},
		)
	}

	return c.Status(fiber.StatusOK).JSON(environment)
}

func CreateEnvironment(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	envName := c.Params("name")
	if err := utils.Validate().Var(envName, "required,alphanum"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment name!"},
		)
	}

	newEnv := models.Environment{Name: envName, UserID: userSessionID}
	var environment models.Environment
	if err := db.Where(&newEnv).First(&environment).Error; err == nil {
		return c.Status(fiber.StatusOK).JSON(
			fiber.Map{"error": fmt.Sprintf(
				"The provided environment '%s' already exists. Please choose a different environment name!", envName),
			},
		)
	}

	// TODO(carlotta): add a limit to how many environments can be created per account

	if err := db.Create(&newEnv).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).SendString(fmt.Sprintf("Successfully created a(n) %s environment!", envName))
}

func DeleteEnvironment(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	id := c.Params("id")
	if err := utils.Validate().Var(id, "required,uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id!"},
		)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		parsedID := utils.MustParseUUID(id)

		var environment models.Environment
		if err := tx.Where(
			&models.Environment{ID: parsedID, UserID: userSessionID},
		).First(&environment).Error; err != nil {
			return c.Status(fiber.StatusOK).JSON(
				fiber.Map{"error": "The provided environment doesn't appear to exist!"},
			)
		}

		var secrets []models.Secret
		if err := tx.Preload("Environments").Not("id", parsedID).Find(
			&secrets, "user_id=?", userSessionID,
		).Error; err != nil || len(secrets) > 0 {
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).SendString(
			fmt.Sprintf("Successfully deleted the %s environment!", environment.Name),
		)
	})

}

type ReqUpdateEnv struct {
	ID          string `json:"id" validate:"required,uuid"`
	UpdatedName string `json:"updatedName" validate:"required"`
}

func UpdateEnvironment(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	data := new(ReqUpdateEnv)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id and updated environment name!"},
		)
	}

	if err := utils.Validate().Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id and updated environment name!"},
		)
	}

	var environment models.Environment
	if err := db.Where(
		&models.Environment{ID: utils.MustParseUUID(data.ID), UserID: userSessionID},
	).First(&environment).Error; err != nil {
		return c.Status(fiber.StatusOK).JSON(
			fiber.Map{"error": "The provided environment doesn't appear to exist."},
		)
	}

	if err := db.Model(&environment).Update("name", &data.UpdatedName).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).SendString(
		fmt.Sprintf("Successfully updated the environment name to %s!", data.UpdatedName),
	)
}
