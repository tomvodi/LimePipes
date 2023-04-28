package apimodel

import "github.com/go-playground/validator/v10"

func NewApimodelValidator() *validator.Validate {
	validate := validator.New()
	validate.SetTagName("binding")
	return validate
}
