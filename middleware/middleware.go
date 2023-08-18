package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mattcarlotta/nvi-api/utils"
)

func RequiresCookieSession(c *fiber.Ctx) error {
	token, err := utils.ValidateSessionToken(c.Cookies("SESSION_TOKEN"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	if len(token.UserId) == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not a valid token."})
	}

	parsedId, err := uuid.Parse(token.UserId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not a valid token."})
	}

	c.Locals("userSessionId", parsedId)

	return c.Next()
}
