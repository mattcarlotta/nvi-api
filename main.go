package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattcarlotta/nvi-api/controllers"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/middleware"
	"github.com/mattcarlotta/nvi-api/utils"
)

func main() {
	database.ConnectDB()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", controllers.Login).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/register", controllers.Register).Methods(http.MethodPost, http.MethodOptions)
	router.Use(middleware.CORS, middleware.Logging)

	authRouter := router.PathPrefix("/").Subrouter()
	authRouter.HandleFunc("/create/environment/{name}", controllers.CreateEnvironment).Methods(http.MethodPost, http.MethodOptions)
	authRouter.HandleFunc("/delete/environment/{name}", controllers.DeleteEnvironment).Methods(http.MethodDelete, http.MethodOptions)
	authRouter.HandleFunc("/update/environment", controllers.UpdateEnvironment).Methods(http.MethodPatch, http.MethodOptions)
	authRouter.Use(middleware.CookieSession)

	var PORT = utils.GetEnv("PORT")
	var API_HOST = utils.GetEnv("API_HOST")
	fmt.Printf("🎧 Listening for incoming requests to %s%s", API_HOST, PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}
