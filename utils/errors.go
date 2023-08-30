package utils

import (
	"fmt"
	"log"
)

type ErrorResponseCode int

const (
	Unknown = iota
	RegisterEmptyBody
	RegisterInvalidBody
	RegisterEmailTaken
	LoginEmptyBody
	LoginInvalidBody
	LoginUnregisteredEmail
	LoginInvalidPassword
	LoginAccountNotVerified
	VerifyAccountInvalidToken
	ResendAccountVerificationInvalidEmail
)

var ErrorCode = map[ErrorResponseCode]string{
	Unknown:                               "E000",
	RegisterEmptyBody:                     "E001",
	RegisterInvalidBody:                   "E002",
	RegisterEmailTaken:                    "E003",
	LoginEmptyBody:                        "E004",
	LoginInvalidBody:                      "E005",
	LoginUnregisteredEmail:                "E006",
	LoginInvalidPassword:                  "E007",
	LoginAccountNotVerified:               "E008",
	VerifyAccountInvalidToken:             "E009",
	ResendAccountVerificationInvalidEmail: "E010",
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
