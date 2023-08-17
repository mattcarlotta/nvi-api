package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Secret struct {
	Environment string `json:"environment"`
	Description string `json:"description"`
	Key         string `json:"key"`
	Content     string `json:"content"`
}

type Secrets []Secret

func allSecrets(c *fiber.Ctx) error {
	secrets := Secrets{
		Secret{
			Environment: "staging",
			Description: "This is an ultra secret key",
			Key:         "BASIC_ENV",
			Content:     "Super secret key",
		},
	}
	fmt.Println("All secrets endpoint")
	return c.Status(fiber.StatusOK).JSON(secrets)
}
