package music_model

import "banduslib/internal/common/music_model/symbols"

type Symbol struct {
	Note *symbols.Note `yaml:"note,omitempty"`
}
