package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/controllers"
	"github.com/mattcarlotta/nvi-api/middlewares"
)

func SecretRoutes(app *fiber.App) {
	secret := app.Group("/")
	secret.Get("/secret/:id", middlewares.RequiresCookieSession, controllers.GetSecretBySecretID)
	secret.Get("/secrets/:id", middlewares.RequiresCookieSession, controllers.GetSecretsByEnvironmentID)
	secret.Post("/create/secret", middlewares.RequiresCookieSession, controllers.CreateSecret)
	secret.Delete("/delete/secret/:id", middlewares.RequiresCookieSession, controllers.DeleteSecret)
	secret.Put("/update/secret/", middlewares.RequiresCookieSession, controllers.UpdateSecret)
}
