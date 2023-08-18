package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
)

type ReqEnv struct {
	ID          string `json:"id"`
	UpdatedName string `json:"updatedName"`
}

func GetAllEnvironments(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := c.Locals("userSessionId").(uuid.UUID)

	var environments []models.Environment
	db.Where(&models.Environment{UserId: userSessionId}).Find(&environments)

	return c.Status(fiber.StatusOK).JSON(environments)
}

func GetEnvironmentById(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := c.Locals("userSessionId").(uuid.UUID)

	id := c.Params("id")
	if len(id) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id!"},
		)
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id!"},
		)
	}

	var environment models.Environment
	if err := db.Where(
		&models.Environment{ID: parsedId, UserId: userSessionId},
	).First(&environment).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "The provided environment doesn't appear to exist!"},
		)
	}

	return c.Status(fiber.StatusOK).JSON(environment)
}

func CreateEnvironment(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := c.Locals("userSessionId").(uuid.UUID)

	envName := c.Params("name")
	if len(envName) == 0 {
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
	userSessionId := c.Locals("userSessionId").(uuid.UUID)

	id := c.Params("id")
	if len(id) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id!"},
		)
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id!"},
		)
	}

	var environment models.Environment
	if err := db.Where(
		&models.Environment{ID: parsedId, UserId: userSessionId},
	).First(&environment).Error; err != nil {
		return c.Status(fiber.StatusOK).JSON(
			fiber.Map{"error": "The provided environment doesn't appear to exist!"},
		)
	}

	if err := db.Delete(&environment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).SendString(
		fmt.Sprintf("Successfully deleted the %s environment!", environment.Name),
	)
}

func UpdateEnvironment(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := c.Locals("userSessionId").(uuid.UUID)

	data := new(ReqEnv)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id and updated environment name!"},
		)
	}

	// TODO(carlotta): Add field validations for "id" and "updatedName"
	if len(data.ID) == 0 || len(data.UpdatedName) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id and updated environment name!"},
		)
	}

	parsedId, err := uuid.Parse(data.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid environment id!"},
		)
	}

	var environment models.Environment
	if err := db.Where(
		&models.Environment{ID: parsedId, UserId: userSessionId},
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
