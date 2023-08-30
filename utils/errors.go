package utils

import (
	"fmt"
	"log"
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
	CreateEnvironmentInvalidName
	CreateEnvironmentNameTaken
	GetEnvironmentInvalidID
	GetEnvironmentNonExistentID
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
	CreateEnvironmentInvalidName:          "E014",
	CreateEnvironmentNameTaken:            "E015",
}

type ResponseError struct {
	Resource string `json:"resource"`
	Error    string `json:"error"`
}

func JSONError(code ErrorResponseCode) ResponseError {
	log.Printf("Error: %s", ErrorCode[code])
	return ResponseError{
		Resource: fmt.Sprintf("https://github.com/mattcarlotta/nvi-api/blob/main/ERRORS.md#%s", ErrorCode[code]),
		Error:    ErrorCode[code],
	}
}

func UnknownJSONError(err error) ResponseError {
	log.Printf("An unknown error occured: %s", err.Error())
	return ResponseError{
		Resource: fmt.Sprintf("https://github.com/mattcarlotta/nvi-api/blob/main/ERRORS.md#%s", ErrorCode[Unknown]),
		Error:    err.Error(),
	}
}
