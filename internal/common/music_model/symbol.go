package music_model

import (
	"banduslib/internal/common/music_model/symbols"
	"banduslib/internal/common/music_model/symbols/time_line"
)

type Symbol struct {
	Note        *symbols.Note       `yaml:"note,omitempty"`
	Rest        *symbols.Rest       `yaml:"rest,omitempty"`
	TimeLine    *time_line.TimeLine `yaml:"timeLine,omitempty"`
	TempoChange uint64              `yaml:"tempoChange,omitempty"`
}

func (s *Symbol) IsNote() bool {
	return s.Note != nil
}

func (s *Symbol) IsValidNote() bool {
	if s.Note != nil {
		return s.Note.IsValid()
	}

	return false
}
