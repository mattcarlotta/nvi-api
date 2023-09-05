package utils

import (
	"fmt"
	"log"
	"os"
)

type ErrorResponseCode int

const (
	Unknown = iota
	RegisterInvalidBody
	RegisterEmailTaken
	LoginInvalidBody
	LoginUnregisteredEmail
	LoginInvalidPassword
	LoginAccountNotVerified
	VerifyAccountInvalidToken
	ResendAccountVerificationInvalidEmail
	SendResetPasswordInvalidEmail
	UpdatePasswordInvalidBody
	UpdatePasswordInvalidToken
	GetEnvironmentInvalidID
	GetEnvironmentNonExistentID
	CreateEnvironmentInvalidBody
	CreateEnvironmentInvalidProjectID
	CreateEnvironmentNameTaken
	DeleteEnvironmentInvalidID
	DeleteEnvironmentNonExistentID
	UpdateEnvironmentInvalidBody
	UpdateEnvironmentNonExistentID
	GetSecretInvalidID
	GetSecretNonExistentID
	GetSecretsByEnvInvalidID
	GetSecretsByEnvNonExistentID
	CreateSecretInvalidBody
	CreateSecretNonExistentProject
	CreateSecretNonExistentEnv
	CreateSecretKeyAlreadyExists
	DeleteSecretInvalidID
	DeleteSecretNonExistentID
	UpdateSecretInvalidBody
	UpdateSecretInvalidID
	UpdateSecretNonExistentEnv
	UpdateSecretKeyAlreadyExists
)

var ErrorCode = map[ErrorResponseCode]string{
	Unknown:                               "E000",
	RegisterInvalidBody:                   "E001",
	RegisterEmailTaken:                    "E002",
	LoginInvalidBody:                      "E003",
	LoginUnregisteredEmail:                "E004",
	LoginInvalidPassword:                  "E005",
	LoginAccountNotVerified:               "E006",
	VerifyAccountInvalidToken:             "E007",
	ResendAccountVerificationInvalidEmail: "E008",
	SendResetPasswordInvalidEmail:         "E009",
	UpdatePasswordInvalidBody:             "E010",
	UpdatePasswordInvalidToken:            "E011",
	GetEnvironmentInvalidID:               "E012",
	GetEnvironmentNonExistentID:           "E013",
	CreateEnvironmentInvalidBody:          "E014",
	CreateEnvironmentInvalidProjectID:     "E015",
	CreateEnvironmentNameTaken:            "E016",
	DeleteEnvironmentInvalidID:            "E017",
	DeleteEnvironmentNonExistentID:        "E018",
	UpdateEnvironmentInvalidBody:          "E019",
	UpdateEnvironmentNonExistentID:        "E020",
	GetSecretInvalidID:                    "E021",
	GetSecretNonExistentID:                "E022",
	GetSecretsByEnvInvalidID:              "E023",
	GetSecretsByEnvNonExistentID:          "E024",
	CreateSecretInvalidBody:               "E025",
	CreateSecretNonExistentProject:        "E026",
	CreateSecretNonExistentEnv:            "E027",
	CreateSecretKeyAlreadyExists:          "E028",
	DeleteSecretInvalidID:                 "E029",
	DeleteSecretNonExistentID:             "E030",
	UpdateSecretInvalidBody:               "E031",
	UpdateSecretInvalidID:                 "E032",
	UpdateSecretNonExistentEnv:            "E033",
	UpdateSecretKeyAlreadyExists:          "E034",
}

type ResponseError struct {
	Resource string `json:"resource"`
	Error    string `json:"error"`
}

func JSONError(code ErrorResponseCode) ResponseError {
	if os.Getenv("IN_TESTING") != "true" {
		log.Printf("Error: %s", ErrorCode[code])
	}
	return ResponseError{
		Resource: fmt.Sprintf("https://github.com/mattcarlotta/nvi-api/blob/main/ERRORS.md#%s", ErrorCode[code]),
		Error:    ErrorCode[code],
	}
}

func UnknownJSONError(err error) ResponseError {
	if os.Getenv("IN_TESTING") != "true" {
		log.Printf("An unknown error occured: %s", err.Error())
	}
	return ResponseError{
		Resource: fmt.Sprintf("https://github.com/mattcarlotta/nvi-api/blob/main/ERRORS.md#%s", ErrorCode[Unknown]),
		Error:    err.Error(),
	}
}
