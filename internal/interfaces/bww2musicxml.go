package interfaces

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/musicxml/model"
)

//go:generate mockgen -source bww2musicxml.go -destination ./mocks/mock_bww2musicxml.go

type Bww2Musicxml interface {
	Convert(music []music_model.Tune) (*model.Score, error)
}
