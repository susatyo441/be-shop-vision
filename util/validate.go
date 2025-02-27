package util

import "github.com/go-playground/validator/v10"

func ValidateStruct(obj any) error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(obj)
}
