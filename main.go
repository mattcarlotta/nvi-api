package main

import (
	"context"
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

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("SESSION_TOKEN")
		if err != nil {
			utils.SendErrorResponse(res, http.StatusUnauthorized, "You must be logged in order to do that!")
			return
		} else if len(cookie.Value) == 0 {
			utils.SendErrorResponse(res, http.StatusUnauthorized, "You must be logged in order to do that!")
			return
		}

		data, err := utils.ValidateSessionToken(cookie.Value)
		if err != nil {
			utils.SendErrorResponse(res, http.StatusUnauthorized, err.Error())
			return
		}
		ctx := context.WithValue(req.Context(), "userSessionId", data.UserId)

		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

func main() {
	database.ConnectDB()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", controllers.Login).Methods(http.MethodPost)
	router.HandleFunc("/register", controllers.CreateUser).Methods(http.MethodPost)
	router.Use(LoggingMiddleware)

	authRouter := router.PathPrefix("/").Subrouter()
	authRouter.HandleFunc("/create/environment/{name}", controllers.CreateEnvironment).Methods(http.MethodPost)
	authRouter.Use(SessionMiddleware)

	var PORT = utils.GetEnv("PORT")
	fmt.Printf("🎧 Listening for requests on port %s", PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}
