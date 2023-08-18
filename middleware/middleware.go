package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mattcarlotta/nvi-api/utils"
)

func RequiresCookieSession(c *fiber.Ctx) error {
	token, err := utils.ValidateSessionToken(c.Cookies("SESSION_TOKEN"))
	if err != nil {
		utils.SetSessionCookie(c, "", time.Unix(0, 0))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	parsedId, err := uuid.Parse(token.UserId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not a valid token."})
	}

	c.Locals("userSessionId", parsedId)

	return c.Next()
}
