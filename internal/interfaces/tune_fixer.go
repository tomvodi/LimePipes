package interfaces

import "github.com/tomvodi/limepipes/internal/common/music_model"

//go:generate mockgen -source tune_fixer.go -destination ./mocks/mock_tune_fixer.go.go

// TuneFixer fixes a tune's meta data like type composer and arranger
type TuneFixer interface {
	Fix(muMo music_model.MusicModel)
}
