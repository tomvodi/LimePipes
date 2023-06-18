package timeline_normalizer

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/symbols/time_line"
)

type MusicChunk struct {
	TimeLineType time_line.TimeLineType

	// StartSymbols are a list of symbols when the music chunk doesn't start with a
	// complete measure
	StartSymbols []*music_model.Symbol

	// Measures contains the full measures of the chunk after StartSymbols and before
	// EndSymbols
	Measures []*music_model.Measure

	// EndSymbols are a list of symbols when the music chunk doesn't end with a
	// complete measure
	EndSymbols []*music_model.Symbol
}
