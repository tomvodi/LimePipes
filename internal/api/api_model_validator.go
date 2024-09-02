package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/tomvodi/limepipes/internal/apigen/apimodel"
)

type ModelValidator struct {
	validator *validator.Validate
}

func (v *ModelValidator) ValidateUpdateTune(tuneUpd apimodel.UpdateTune) error {
	return v.validator.Struct(tuneUpd)
}

func (v *ModelValidator) ValidateUpdateSet(setUpd apimodel.UpdateSet) error {
	return v.validator.Struct(setUpd)
}

func NewAPIModelValidator(
	validator *validator.Validate,
) *ModelValidator {
	return &ModelValidator{
		validator: validator,
	}
}
