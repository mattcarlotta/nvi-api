package utils

import "github.com/gofiber/fiber/v2"

func SendErrorResponse(c *fiber.Ctx, code int, message string) error {
	if len(message) == 0 {
		return c.Status(code).Send(nil)
	}
	return c.Status(code).JSON(fiber.Map{"error": message})
}
