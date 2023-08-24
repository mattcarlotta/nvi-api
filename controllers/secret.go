package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
)

type ReqSecret struct {
	EnvironmentIds []string `json:"environmentIds"`
	Name           string   `json:"name"`
	Content        string   `json:"content"`
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
	if err := db.Preload("Environments").Find(
		&secrets, "name=? AND user_id=?", data.Name, userSessionId,
	).Error; err == nil {
		for _, secret := range secrets {
			for _, env := range secret.Environments {
				for _, id := range environmentIds {
					if env.ID == id {
						return c.Status(fiber.StatusConflict).JSON(
							fiber.Map{"error": fmt.Sprintf(
								"The key '%s' already exists for the '%s' environment.", data.Name, env.Name),
							},
						)
					}
				}
			}
		}
	}

	newSecret := models.Secret{
		Name:         data.Name,
		Content:      data.Content,
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
