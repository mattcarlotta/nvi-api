package utils

import (
	"log"
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func validateUUIDArray(fl validator.FieldLevel) bool {
	arr := fl.Field()
	if arr.Kind() != reflect.Slice || arr.Len() == 0 {
		return false
	}

	for i := 0; i < arr.Len(); i++ {
		item := arr.Index(i).Interface()

		str, ok := item.(string)
		if !ok {
			return false
		}

		_, err := ParseUUID(str)
		if err != nil {
			return false
		}
	}

	return true
}

var environmentNameRegex = regexp.MustCompile("^[a-zA-Z0-9_]+$")

func validateEnvironmentName(fl validator.FieldLevel) bool {
	return environmentNameRegex.MatchString(fl.Field().String())
}

func Validate() *validator.Validate {
	if validate == nil {
		validate = validator.New()
		if err := validate.RegisterValidation("uuidarray", validateUUIDArray); err != nil {
			log.Fatalf("Unable to register uuidarray validator: %s", err.Error())
		}
		if err := validate.RegisterValidation("envname", validateEnvironmentName); err != nil {
			log.Fatalf("Unable to register envname validator: %s", err.Error())
		}
		return validate
	}
	return validate
}
