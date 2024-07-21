package interfaces

import (
	"github.com/tomvodi/limepipes/internal/common/music_model"
	"github.com/tomvodi/limepipes/internal/musicxml/model"
)

type MusicxmlExporter interface {
	Export(musicModel music_model.MusicModel) (*model.Score, error)
}
