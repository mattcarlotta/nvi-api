package utils

import (
	"reflect"

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

func Validate() *validator.Validate {
	if validate == nil {
		validate = validator.New()
		validate.RegisterValidation("uuidarray", validateUUIDArray)
		return validate
	} else {
		return validate
	}
}
