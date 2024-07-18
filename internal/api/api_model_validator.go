package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/tomvodi/limepipes/internal/api_gen/apimodel"
	"github.com/tomvodi/limepipes/internal/interfaces"
)

type apiModelValidator struct {
	validator *validator.Validate
}

func (v *apiModelValidator) ValidateUpdateTune(tuneUpd apimodel.UpdateTune) error {
	return v.validator.Struct(tuneUpd)
}

func (v *apiModelValidator) ValidateUpdateSet(setUpd apimodel.UpdateSet) error {
	return v.validator.Struct(setUpd)
}

func NewApiModelValidator(
	validator *validator.Validate,
) interfaces.ApiModelValidator {
	return &apiModelValidator{
		validator: validator,
	}
}
