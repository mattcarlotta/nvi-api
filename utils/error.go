package utils

import "github.com/gofiber/fiber/v2"

func SendErrorResponse(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{"error": message})
}
