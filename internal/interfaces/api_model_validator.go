package interfaces

import "github.com/tomvodi/limepipes/internal/api_gen/apimodel"

type ApiModelValidator interface {
	ValidateUpdateTune(tuneUpd apimodel.UpdateTune) error
	ValidateUpdateSet(tuneUpd apimodel.UpdateSet) error
}
