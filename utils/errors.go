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
	GetAllEnvironmentsInvalidProjectID
	GetAllEnvironmentsNonExistentID
	GetEnvironmentInvalidID
	GetEnvironmentInvalidName
	GetEnvironmentInvalidProjectID
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
	SearchForSecretsByEnvAndSecretInvalidKey
	CreateProjectOverLimit
)

var ErrorCode = map[ErrorResponseCode]string{
	Unknown:                                  "E000",
	RegisterInvalidBody:                      "E001",
	RegisterEmailTaken:                       "E002",
	LoginInvalidBody:                         "E003",
	LoginUnregisteredEmail:                   "E004",
	LoginInvalidPassword:                     "E005",
	LoginAccountNotVerified:                  "E006",
	VerifyAccountInvalidToken:                "E007",
	ResendAccountVerificationInvalidEmail:    "E008",
	SendResetPasswordInvalidEmail:            "E009",
	UpdatePasswordInvalidBody:                "E010",
	UpdatePasswordInvalidToken:               "E011",
	GetAllEnvironmentsInvalidProjectID:       "E012",
	GetAllEnvironmentsNonExistentID:          "E013",
	GetEnvironmentInvalidID:                  "E012",
	GetEnvironmentNonExistentID:              "E013",
	GetEnvironmentInvalidName:                "E014",
	GetEnvironmentInvalidProjectID:           "E015",
	GetEnvironmentNonExistentName:            "E016",
	CreateEnvironmentInvalidBody:             "E017",
	CreateEnvironmentInvalidProjectID:        "E018",
	CreateEnvironmentNameTaken:               "E019",
	DeleteEnvironmentInvalidID:               "E020",
	DeleteEnvironmentNonExistentID:           "E021",
	UpdateEnvironmentInvalidBody:             "E022",
	UpdateEnvironmentInvalidProjectID:        "E023",
	UpdateEnvironmentNonExistentID:           "E024",
	UpdateEnvironmentNameTaken:               "E025",
	GetSecretInvalidID:                       "E026",
	GetSecretNonExistentID:                   "E027",
	GetSecretsByEnvInvalidID:                 "E028",
	GetSecretsByEnvNonExistentID:             "E029",
	CreateSecretInvalidBody:                  "E030",
	CreateSecretNonExistentProject:           "E031",
	CreateSecretNonExistentEnv:               "E032",
	CreateSecretKeyAlreadyExists:             "E033",
	DeleteSecretInvalidID:                    "E034",
	DeleteSecretNonExistentID:                "E035",
	UpdateSecretInvalidBody:                  "E036",
	UpdateSecretInvalidID:                    "E037",
	UpdateSecretNonExistentEnv:               "E038",
	UpdateSecretKeyAlreadyExists:             "E039",
	GetProjectInvalidID:                      "E040",
	GetProjectNonExistentID:                  "E041",
	GetProjectInvalidName:                    "E042",
	GetProjectNonExistentName:                "E043",
	CreateProjectInvalidName:                 "E044",
	CreateProjectNameTaken:                   "E045",
	DeleteProjectInvalidID:                   "E046",
	DeleteProjectNonExistentID:               "E047",
	UpdateProjectInvalidBody:                 "E048",
	UpdateProjectNonExistentID:               "E049",
	UpdateProjectNameTaken:                   "E050",
	SearchForSecretsByEnvAndSecretInvalidKey: "E051",
	CreateProjectOverLimit:                   "E052",
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
