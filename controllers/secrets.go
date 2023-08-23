package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	// "github.com/mattcarlotta/nvi-api/utils"
)

type ReqSecret struct {
	Environment string `json:"environment"`
	Name        string `json:"name"`
	Content     string `json:"content"`
}

func CreateSecret(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := c.Locals("userSessionId").(uuid.UUID)

	data := new(ReqSecret)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide valid environments and a secret key name with content!"},
		)
	}

	// TODO(carlotta): Add field validations for "environment," "name," and "content"
	if len(data.Environment) == 0 || len(data.Name) == 0 || len(data.Content) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide valid environments and a secret key name with content!"},
		)
	}

	parsedEnvId, err := uuid.Parse(data.Environment)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id!"},
		)
	}

	var environment models.Environment
	if err := db.Where(&models.Environment{ID: parsedEnvId, UserId: userSessionId}).First(&environment).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "The provided environment doesn't appear to exist!"},
		)
	}

	newSecret := models.Secret{Name: data.Name, UserId: userSessionId, EnvironmentId: parsedEnvId}
	var secret models.Secret
	if err := db.Where(&newSecret).First(&secret).Error; err == nil {
		return c.Status(fiber.StatusOK).JSON(
			fiber.Map{"error": fmt.Sprintf(
				"The provided name '%s' already exists for the selected environment.", data.Name),
			},
		)
	}

	newSecret.Content = data.Content
	if err := db.Create(&newSecret).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).SendString(
		fmt.Sprintf("Successfully created %s!", data.Name),
	)
}
