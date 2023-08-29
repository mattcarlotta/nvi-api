package utils

type ErrorResponseCode int

const (
	RegisterEmptyBody = iota
	RegisterInvalidBody
	RegisterEmailTaken
)

var ErrorCode = map[ErrorResponseCode]string{
	RegisterEmptyBody:   "E001",
	RegisterInvalidBody: "E002",
	RegisterEmailTaken:  "E003",
}
