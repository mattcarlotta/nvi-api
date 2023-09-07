package utils

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetSessionCookie(c *fiber.Ctx, value string, expires time.Time) {
	cookie := fiber.Cookie{
		Name:     "SESSION_TOKEN",
		Value:    value,
		Expires:  expires,
		Path:     "/",
		HTTPOnly: false,
		Secure:   os.Getenv("IN_PRODUCTION") == "true",
		SameSite: "Lax",
	}

	c.Cookie(&cookie)
}
