package interfaces

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/musicxml/model"
)

//go:generate mockgen -source musicxml_exporter.go -destination ./mocks/mock_musicxml_exporter.go

type MusicxmlExporter interface {
	Export(musicModel music_model.MusicModel) (*model.Score, error)
}
