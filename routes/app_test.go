package routes

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"testing"

	"github.com/mattcarlotta/nvi-api/database"
)

var app *fiber.App

func TestMain(m *testing.M) {
	database.CreateConnection()

	app = fiber.New()

	UserRoutes(app)
	EnvironmentRoutes(app)
	SecretRoutes(app)

	// log.Fatal(app.Listen(utils.GetEnv("PORT")))
	os.Exit(m.Run())
}
