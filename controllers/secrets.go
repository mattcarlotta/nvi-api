package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Secret struct {
	Environment string `json:"environment"`
	Description string `json:"description"`
	Key         string `json:"key"`
	Content     string `json:"content"`
}

type Secrets []Secret

func allSecrets(res http.ResponseWriter, req *http.Request) {
	secrets := Secrets{
		Secret{
			Environment: "staging",
			Description: "This is an ultra secret key",
			Key:         "BASIC_ENV",
			Content:     "Super secret key",
		},
	}
	fmt.Println("All secrets endpoint")
	json.NewEncoder(res).Encode(secrets)
}
