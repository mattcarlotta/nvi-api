package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mattcarlotta/nvi-api/controllers"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/utils"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

func main() {
	database.ConnectDB()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", controllers.Login).Methods(http.MethodPost)
	router.HandleFunc("/create/user", controllers.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/secrets", allSecrets).Methods(http.MethodGet)
	router.Use(LoggingMiddleware)

	var PORT = utils.GetEnv("PORT")
	fmt.Printf("🎧 Listening for requests on port %s", PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}
