package routes

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
)

var app *fiber.App

func sendAppRequest(req *http.Request) *http.Response {
	res, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to that controller")
	}

	return res
}

func TestMain(m *testing.M) {
	db := database.CreateConnection()

	if err := db.Migrator().DropTable(&models.User{}); err != nil {
		log.Fatalf("Unable to drop user table: %s", err.Error())
	}
	if err := db.Migrator().DropTable(&models.Environment{}); err != nil {
		log.Fatalf("Unable to drop environment table: %s", err.Error())
	}
	if err := db.Migrator().DropTable(&models.Secret{}); err != nil {
		log.Fatalf("Unable to drop secret table: %s", err.Error())
	}

	if err := db.AutoMigrate(&models.User{}, &models.Environment{}, &models.Secret{}); err != nil {
		log.Fatalf("Unable to migrate models: %s", err.Error())
	}

	app = fiber.New()

	CLIRoutes(app)
	UserRoutes(app)
	EnvironmentRoutes(app)
	SecretRoutes(app)
	ProjectRoutes(app)

	os.Exit(m.Run())
}
