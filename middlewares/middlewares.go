package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
	"github.com/mattcarlotta/nvi-api/utils"
)

func Setup(app *fiber.App) {
	app.Use(cors.New(
		cors.Config{
			AllowOrigins:     utils.GetEnv("CLIENT_HOST"),
			AllowHeaders:     "Origin, Content-Type, Accept",
			AllowCredentials: true,
		},
	))
	app.Use(helmet.New())
	app.Use(encryptcookie.New(
		encryptcookie.Config{
			Key: utils.GetEnv("COOKIE_KEY"),
		}),
	)
	app.Use(compress.New(
		compress.Config{
			Level: compress.LevelBestSpeed,
		}),
	)
	app.Use(logger.New())
}

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
