package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/controllers"
	"github.com/mattcarlotta/nvi-api/middlewares"
)

func CLIRoutes(app *fiber.App) {
	cli := app.Group("/")
	cli.Get("/cli/secrets", middlewares.RequiresAPIKey, controllers.GetSecretsByAPIKey)
	cli.Get("/cli/projects", middlewares.RequiresAPIKey, controllers.GetProjectsByAPIKey)
	cli.Get("/cli/environments", middlewares.RequiresAPIKey, controllers.GetEnvironmentsByAPIKey)
}
