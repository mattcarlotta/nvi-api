package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/controllers"
	"github.com/mattcarlotta/nvi-api/middlewares"
)

func ProjectRoutes(app *fiber.App) {
	project := app.Group("/")
	project.Get("/project/id/:id", middlewares.RequiresCookieSession, controllers.GetProjectByID)
	project.Get("/project/name/:name", middlewares.RequiresCookieSession, controllers.GetProjectByName)
	project.Get("/projects", middlewares.RequiresCookieSession, controllers.GetAllProjects)
	project.Post("/create/project/:name", middlewares.RequiresCookieSession, controllers.CreateProject)
	project.Delete("/delete/project/:id", middlewares.RequiresCookieSession, controllers.DeleteProject)
	project.Put("/update/project", middlewares.RequiresCookieSession, controllers.UpdateProject)
}
