package interfaces

import "github.com/tomvodi/limepipes/internal/apigen/apimodel"

type APIModelValidator interface {
	ValidateUpdateTune(tuneUpd apimodel.UpdateTune) error
	ValidateUpdateSet(tuneUpd apimodel.UpdateSet) error
}
