package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/controllers"
	"github.com/mattcarlotta/nvi-api/middlewares"
)

func Setup(app *fiber.App) {
	user := app.Group("/")
	user.Post("/register", controllers.Register)
	user.Post("/login", controllers.Login)
	user.Post("/logout", controllers.Logout)
	user.Patch("/verify/account", controllers.VerifyAccount)
	user.Post("/reverify/account", controllers.ResendAccountVerification)
	user.Post("/reset/password", controllers.SendResetPasswordEmail)
	user.Patch("/update/password", controllers.UpdatePassword)
	user.Get("/account", middlewares.RequiresCookieSession, controllers.GetAccountInfo)
	user.Delete("/delete/account", middlewares.RequiresCookieSession, controllers.DeleteAccount)

	environment := app.Group("/")
	environment.Get("/environments", middlewares.RequiresCookieSession, controllers.GetAllEnvironments)
	environment.Get("/environment/:id", middlewares.RequiresCookieSession, controllers.GetEnvironmentById)
	environment.Post("/create/environment/:name", middlewares.RequiresCookieSession, controllers.CreateEnvironment)
	environment.Delete("/delete/environment/:id", middlewares.RequiresCookieSession, controllers.DeleteEnvironment)
	environment.Patch("/update/environment", middlewares.RequiresCookieSession, controllers.UpdateEnvironment)

	secret := app.Group("/")
	secret.Post("/create/secret", middlewares.RequiresCookieSession, controllers.CreateSecret)
}
