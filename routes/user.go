package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/controllers"
	"github.com/mattcarlotta/nvi-api/middlewares"
)

func UserRoutes(app *fiber.App) {
	user := app.Group("/")
	user.Post("/register", controllers.Register)
	user.Post("/login", controllers.Login)
	user.Get("/loggedin", middlewares.RequiresCookieSession, controllers.Loggedin)
	user.Post("/logout", controllers.Logout)
	user.Patch("/verify/account", controllers.VerifyAccount)
	user.Patch("/reverify/account", controllers.ResendAccountVerification)
	user.Patch("/reset/password", controllers.SendResetPasswordEmail)
	user.Patch("/update/password", controllers.UpdatePassword)
	user.Patch("/update/apikey", middlewares.RequiresCookieSession, controllers.UpdateAPIKey)
	user.Get("/account", middlewares.RequiresCookieSession, controllers.GetAccountInfo)
	user.Delete("/delete/account", middlewares.RequiresCookieSession, controllers.DeleteAccount)
}
