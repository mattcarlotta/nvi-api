package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/utils"
)

func RequiresCookieSession() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies("SESSION_TOKEN")
		if len(cookie) == 0 {
			return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "You must be logged in order to do that!")
		}

		token, err := utils.ValidateSessionToken(cookie)
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusUnauthorized, err.Error())
		}

		c.Locals("userSessionId", token.UserId)

		return c.Next()
	}
}
