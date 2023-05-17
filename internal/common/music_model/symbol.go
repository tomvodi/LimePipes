package music_model

import (
	"banduslib/internal/common/music_model/symbols"
	"banduslib/internal/common/music_model/symbols/tuplet"
)

type Symbol struct {
	Note   *symbols.Note  `yaml:"note,omitempty"`
	Rest   *symbols.Rest  `yaml:"rest,omitempty"`
	Tuplet *tuplet.Tuplet `yaml:"tuplet,omitempty"`
}
