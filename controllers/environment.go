package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
	"gorm.io/gorm"
)

type ReqEnv struct {
	ID          string `json:"id" validate:"required,uuid"`
	UpdatedName string `json:"updatedName" validate:"required"`
}

func GetAllEnvironments(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := utils.GetSessionId(c)

	var environments []models.Environment
	db.Where(&models.Environment{UserId: userSessionId}).Find(&environments)

	return c.Status(fiber.StatusOK).JSON(environments)
}

func GetEnvironmentById(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := utils.GetSessionId(c)

	id := c.Params("id")
	err := utils.Validate().Var(id, "required,uuid")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id!"},
		)
	}

	var environment models.Environment
	if err := db.Where(
		&models.Environment{ID: utils.MustParseUUID(id), UserId: userSessionId},
	).First(&environment).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "The provided environment doesn't appear to exist!"},
		)
	}

	return c.Status(fiber.StatusOK).JSON(environment)
}

func CreateEnvironment(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := utils.GetSessionId(c)

	envName := c.Params("name")
	err := utils.Validate().Var(envName, "required,alphanum")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment name!"},
		)
	}

	newEnv := models.Environment{Name: envName, UserId: userSessionId}
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
	userSessionId := utils.GetSessionId(c)

	id := c.Params("id")
	err := utils.Validate().Var(id, "required,uuid")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id!"},
		)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		parsedId := utils.MustParseUUID(id)

		var environment models.Environment
		if err := tx.Where(
			&models.Environment{ID: parsedId, UserId: userSessionId},
		).First(&environment).Error; err != nil {
			return c.Status(fiber.StatusOK).JSON(
				fiber.Map{"error": "The provided environment doesn't appear to exist!"},
			)
		}

		var secrets []models.Secret
		if err := tx.Preload("Environments").Not("id", parsedId).Find(
			&secrets, "user_id=?", userSessionId,
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

func UpdateEnvironment(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := utils.GetSessionId(c)

	data := new(ReqEnv)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id and updated environment name!"},
		)
	}

	err := utils.Validate().Struct(data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id and updated environment name!"},
		)
	}

	var environment models.Environment
	if err := db.Where(
		&models.Environment{ID: utils.MustParseUUID(data.ID), UserId: userSessionId},
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
