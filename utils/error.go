package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendErrorResponse(res http.ResponseWriter, code int, message string) {
	response, err := json.Marshal(ErrorResponse{Error: message})
	if err != nil {
		log.Fatalf("Unable to format json error message: %v", err)
	}

	fmt.Printf("Error: %s", message)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	res.Write(response)
}
