package utils

type ErrorResponse struct {
	Error string `json:"error"`
}

type ErrorResponseCode int

const (
	RegisterEmptyBody = iota
	RegisterInvalidBody
	RegisterEmailTaken
	LoginEmptyBody
	LoginInvalidBody
	LoginUnregisteredEmail
	LoginInvalidPassword
	LoginAccountNotVerified
)

var ErrorCode = map[ErrorResponseCode]string{
	RegisterEmptyBody:       "E001",
	RegisterInvalidBody:     "E002",
	RegisterEmailTaken:      "E003",
	LoginEmptyBody:          "E004",
	LoginInvalidBody:        "E005",
	LoginUnregisteredEmail:  "E006",
	LoginInvalidPassword:    "E007",
	LoginAccountNotVerified: "E008",
}
