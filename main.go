package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mattcarlotta/nvi-api/controllers"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/middleware"
	"github.com/mattcarlotta/nvi-api/utils"
)

func main() {
	database.CreateConnection()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     utils.GetEnv("CLIENT_HOST"),
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))
	app.Use(logger.New())

	user := app.Group("/")
	user.Post("/register", controllers.Register)
	user.Post("/login", controllers.Login)
	user.Post("/logout", controllers.Logout)
	user.Patch("/verify/account", controllers.VerifyAccount)
	user.Post("/reverify/account", controllers.ResendAccountVerification)
	user.Post("/reset/password", controllers.SendResetPasswordEmail)
	user.Patch("/update/password", controllers.UpdatePassword)
	user.Get("/account", middleware.RequiresCookieSession(), controllers.GetAccountInfo)
	user.Delete("/delete/account", middleware.RequiresCookieSession(), controllers.DeleteAccount)

	environment := app.Group("/")
	environment.Get("/environments", middleware.RequiresCookieSession(), controllers.GetAllEnvironments)
	environment.Get("/environment/:id", middleware.RequiresCookieSession(), controllers.GetEnvironmentById)
	environment.Post("/create/environment/:name", middleware.RequiresCookieSession(), controllers.CreateEnvironment)
	environment.Delete("/delete/environment/:id", middleware.RequiresCookieSession(), controllers.DeleteEnvironment)
	environment.Patch("/update/environment", middleware.RequiresCookieSession(), controllers.UpdateEnvironment)

	var PORT = utils.GetEnv("PORT")
	var API_HOST = utils.GetEnv("API_HOST")
	fmt.Printf("🎧 Listening for incoming requests to %s%s", API_HOST, PORT)
	log.Fatal(app.Listen(PORT))
}
