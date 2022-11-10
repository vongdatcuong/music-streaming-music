package validator_utils

import "github.com/go-playground/validator/v10"

func ValidateStruct(obj any) error {
	return validator.New().Struct(obj)
}
