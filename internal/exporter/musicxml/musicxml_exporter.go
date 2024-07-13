package musicxml

import (
	"github.com/tomvodi/limepipes/internal/common/music_model"
	"github.com/tomvodi/limepipes/internal/interfaces"
	"github.com/tomvodi/limepipes/internal/musicxml/model"
)

type xmlExport struct {
}

func (x *xmlExport) Export(musicModel music_model.MusicModel) (*model.Score, error) {
	//TODO implement me
	panic("implement me")
}

func NewBww2MusicxmlConverter() interfaces.MusicxmlExporter {
	return &xmlExport{}
}
