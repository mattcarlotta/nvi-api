package routes

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
)

var app *fiber.App

func TestMain(m *testing.M) {
	database.CreateConnection()

	app = fiber.New()

	UserRoutes(app)
	EnvironmentRoutes(app)
	SecretRoutes(app)

	os.Exit(m.Run())
}
