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
	GetEnvironmentInvalidName
	GetEnvironmentNonExistentID
	GetEnvironmentNonExistentName
	CreateEnvironmentInvalidBody
	CreateEnvironmentInvalidProjectID
	CreateEnvironmentNameTaken
	DeleteEnvironmentInvalidID
	DeleteEnvironmentNonExistentID
	UpdateEnvironmentInvalidBody
	UpdateEnvironmentInvalidProjectID
	UpdateEnvironmentNameTaken
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
	GetProjectInvalidID
	GetProjectNonExistentID
	GetProjectInvalidName
	GetProjectNonExistentName
	CreateProjectInvalidName
	CreateProjectNameTaken
	DeleteProjectInvalidID
	DeleteProjectNonExistentID
	UpdateProjectInvalidBody
	UpdateProjectNonExistentID
	UpdateProjectNameTaken
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
	GetEnvironmentInvalidName:             "E014",
	GetEnvironmentNonExistentName:         "E015",
	CreateEnvironmentInvalidBody:          "E016",
	CreateEnvironmentInvalidProjectID:     "E017",
	CreateEnvironmentNameTaken:            "E018",
	DeleteEnvironmentInvalidID:            "E019",
	DeleteEnvironmentNonExistentID:        "E020",
	UpdateEnvironmentInvalidBody:          "E021",
	UpdateEnvironmentInvalidProjectID:     "E022",
	UpdateEnvironmentNonExistentID:        "E023",
	UpdateEnvironmentNameTaken:            "E024",
	GetSecretInvalidID:                    "E025",
	GetSecretNonExistentID:                "E026",
	GetSecretsByEnvInvalidID:              "E027",
	GetSecretsByEnvNonExistentID:          "E028",
	CreateSecretInvalidBody:               "E029",
	CreateSecretNonExistentProject:        "E030",
	CreateSecretNonExistentEnv:            "E031",
	CreateSecretKeyAlreadyExists:          "E032",
	DeleteSecretInvalidID:                 "E033",
	DeleteSecretNonExistentID:             "E034",
	UpdateSecretInvalidBody:               "E035",
	UpdateSecretInvalidID:                 "E036",
	UpdateSecretNonExistentEnv:            "E037",
	UpdateSecretKeyAlreadyExists:          "E038",
	GetProjectInvalidID:                   "E039",
	GetProjectNonExistentID:               "E040",
	GetProjectInvalidName:                 "E041",
	GetProjectNonExistentName:             "E042",
	CreateProjectInvalidName:              "E043",
	CreateProjectNameTaken:                "E044",
	DeleteProjectInvalidID:                "E045",
	DeleteProjectNonExistentID:            "E046",
	UpdateProjectInvalidBody:              "E047",
	UpdateProjectNonExistentID:            "E048",
	UpdateProjectNameTaken:                "E049",
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
