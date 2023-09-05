package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/middlewares"
	"github.com/mattcarlotta/nvi-api/routes"
	"github.com/mattcarlotta/nvi-api/utils"
)

func main() {
	database.CreateConnection()

	app := fiber.New(fiber.Config{
		ServerHeader: "nvi-api",
		AppName:      "Nvi API v0.0.1",
	})

	middlewares.Setup(app)

	routes.UserRoutes(app)
	routes.EnvironmentRoutes(app)
	routes.SecretRoutes(app)
	routes.ProjectRoutes(app)

	log.Fatal(app.Listen(utils.GetEnv("PORT")))
}
