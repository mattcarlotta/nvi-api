package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mattcarlotta/nvi-api/utils"
)

func Setup(app *fiber.App) {
	app.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     utils.GetEnv("CLIENT_HOST"),
				AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
				AllowCredentials: true,
			},
		),
		helmet.New(),
		encryptcookie.New(
			encryptcookie.Config{
				Key: utils.GetEnv("COOKIE_KEY"),
			},
		),
		compress.New(
			compress.Config{
				Level: compress.LevelBestSpeed,
			},
		),
		logger.New(),
	)
}

func RequiresAPIKey(c *fiber.Ctx) error {
	apiKey := c.Query("apiKey")
	if err := utils.Validate().Var(apiKey, "required,alphanum"); err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(
			"a valid apiKey must be supplied in order to use the cli endpoint",
		)
	}

	c.Locals("apiKey", apiKey)

	return c.Next()
}

func RequiresCookieSession(c *fiber.Ctx) error {
	token, err := utils.ValidateSessionToken(c.Cookies("SESSION_TOKEN"))
	if err != nil {
		utils.SetSessionCookie(c, "", time.Unix(0, 0))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	parsedID, err := utils.ParseUUID(token.UserID)
	if err != nil {
		utils.SetSessionCookie(c, "", time.Unix(0, 0))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Not a valid token."})
	}

	c.Locals("userSessionID", parsedID)

	return c.Next()
}
