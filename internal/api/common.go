package api

import "github.com/go-playground/validator/v10"

func NewGinValidator() *validator.Validate {
	validate := validator.New()
	validate.SetTagName("binding")
	return validate
}
