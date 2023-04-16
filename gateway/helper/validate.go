package helper

import (
	"errors"
	"log"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(user interface{}) error {
	val := reflect.ValueOf(user)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		log.Fatalln("the parameter isn't a struct")
	}

	errorValidation := "validation error on field:"
	err := validate.Struct(user)
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
