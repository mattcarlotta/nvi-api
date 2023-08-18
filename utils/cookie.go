package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetSessionCookie(c *fiber.Ctx, value string, expires time.Time) {
	cookie := fiber.Cookie{
		Name:     "SESSION_TOKEN",
		Value:    value,
		Expires:  expires,
		Path:     "/",
		HTTPOnly: true,
		//Secure: true,
	}

	c.Cookie(&cookie)
}
