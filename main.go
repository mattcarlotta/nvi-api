package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func homePage(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "Home page endpoint\n")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/users", allUsers).Methods(http.MethodGet)
	router.HandleFunc("/secrets", allSecrets).Methods(http.MethodGet)

	PORT, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("The ENV 'PORT' must be defined!")
	}

	log.Printf("Listening for requests on port %s", PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}
