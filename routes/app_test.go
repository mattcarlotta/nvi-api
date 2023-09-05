package routes

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
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
	database.CreateConnection()

	app = fiber.New()

	UserRoutes(app)
	EnvironmentRoutes(app)
	SecretRoutes(app)
	ProjectRoutes(app)

	os.Exit(m.Run())
}
