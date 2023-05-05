package interfaces

import (
	"banduslib/internal/common/music_model"
)

//go:generate mockgen -source bww_parser.go -destination ./mocks/mock_bww_parser.go

type BwwParser interface {
	ParseBwwData(data []byte) ([]*music_model.Tune, error)
}
