package main

import (
	"fmt"
	"log"
	// "net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mattcarlotta/nvi-api/controllers"
	"github.com/mattcarlotta/nvi-api/database"
	// "github.com/mattcarlotta/nvi-api/middleware"
	"github.com/mattcarlotta/nvi-api/utils"
)

func main() {
	database.CreateConnection()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: utils.GetEnv("CLIENT_HOST"),
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	api := app.Group("/", logger.New())

	user := api.Group("")
	user.Post("/login", controllers.Login)
	user.Post("/register", controllers.Register)
	user.Post("/logout", controllers.Logout)
	// user.Post("/verify/account", controllers.VerifyAccount)
	// user.Post("/reverify/account", controllers.ResendAccountVerification)
	// user.Post("/reset/password", controllers.SendResetPasswordEmail)
	// user.Put("/update/password", controllers.UpdatePassword)
	// router.Use(middleware.CORS, middleware.Logging)

	// authRouter := router.PathPrefix("/").Subrouter()
	// authRouter.HandleFunc("/delete/account", controllers.DeleteAccount).Methods(http.MethodDelete, http.MethodOptions)
	// authRouter.HandleFunc("/create/environment/{name}", controllers.CreateEnvironment).Methods(http.MethodPost, http.MethodOptions)
	// authRouter.HandleFunc("/delete/environment/{name}", controllers.DeleteEnvironment).Methods(http.MethodDelete, http.MethodOptions)
	// authRouter.HandleFunc("/update/environment", controllers.UpdateEnvironment).Methods(http.MethodPatch, http.MethodOptions)
	// authRouter.Use(middleware.CookieSession)

	var PORT = utils.GetEnv("PORT")
	var API_HOST = utils.GetEnv("API_HOST")
	fmt.Printf("🎧 Listening for incoming requests to %s%s", API_HOST, PORT)
	log.Fatal(app.Listen(PORT))
}
