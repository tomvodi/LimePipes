package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/musicmodel"
	"github.com/tomvodi/limepipes/internal/musicxml/model"
)

type MusicxmlExporter interface {
	Export(musicModel musicmodel.MusicModel) (*model.Score, error)
}
