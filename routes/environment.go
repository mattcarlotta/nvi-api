package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/controllers"
	"github.com/mattcarlotta/nvi-api/middlewares"
)

func EnvironmentRoutes(app *fiber.App) {
	environment := app.Group("/")
	environment.Get("/environment/id/:id", middlewares.RequiresCookieSession, controllers.GetEnvironmentByID)
	environment.Get("/environments/project/:name", middlewares.RequiresCookieSession, controllers.GetAllEnvironmentByProjectName)
	environment.Get("/environment/name", middlewares.RequiresCookieSession, controllers.GetEnvironmentByNameAndProjectID)
	environment.Get("/environments/search", middlewares.RequiresCookieSession, controllers.SearchForEnvironmentsByNameAndProjectID)
	environment.Post("/create/environment", middlewares.RequiresCookieSession, controllers.CreateEnvironment)
	environment.Delete("/delete/environment/:id", middlewares.RequiresCookieSession, controllers.DeleteEnvironment)
	environment.Put("/update/environment", middlewares.RequiresCookieSession, controllers.UpdateEnvironment)
}
