package helper

import (
	"errors"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(entity interface{}) error {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}

	errorValidation := "validation error on field:"
	err := validate.Struct(entity)
	if err != nil {
		for i, err := range err.(validator.ValidationErrors) {
			if i > 0 {
				errorValidation += " & "
			}
			errorValidation += err.StructNamespace()
		}
	} else {
		return nil
	}
	return errors.New(errorValidation)
}
