package musicxml

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/interfaces"
	"banduslib/internal/musicxml/model"
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
