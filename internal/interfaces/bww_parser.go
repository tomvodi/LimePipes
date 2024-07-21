package interfaces

import (
	"github.com/tomvodi/limepipes/internal/common/music_model"
)

type BwwParser interface {
	ParseBwwData(data []byte) (music_model.MusicModel, error)
}
