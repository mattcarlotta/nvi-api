package utils

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func Validate() *validator.Validate {
	if validate == nil {
		validate = validator.New()
		return validate
	} else {
		return validate
	}
}
