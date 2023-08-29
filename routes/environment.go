package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/controllers"
	"github.com/mattcarlotta/nvi-api/middlewares"
)

func EnvironmentRoutes(app *fiber.App) {
	environment := app.Group("/")
	environment.Get("/environment/:id", middlewares.RequiresCookieSession, controllers.GetEnvironmentByID)
	environment.Get("/environments", middlewares.RequiresCookieSession, controllers.GetAllEnvironments)
	environment.Post("/create/environment/:name", middlewares.RequiresCookieSession, controllers.CreateEnvironment)
	environment.Delete("/delete/environment/:id", middlewares.RequiresCookieSession, controllers.DeleteEnvironment)
	environment.Patch("/update/environment", middlewares.RequiresCookieSession, controllers.UpdateEnvironment)
}
