package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendErrorResponse(res http.ResponseWriter, code int, message string) {
	response, _ := json.Marshal(map[string]string{"error": message})

	fmt.Printf("Error: %s", message)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	res.Write(response)
}
