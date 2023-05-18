package music_model

import (
	"banduslib/internal/common/music_model/symbols"
	"banduslib/internal/common/music_model/symbols/time_line"
	"banduslib/internal/common/music_model/symbols/tuplet"
)

type Symbol struct {
	Note     *symbols.Note       `yaml:"note,omitempty"`
	Rest     *symbols.Rest       `yaml:"rest,omitempty"`
	Tuplet   *tuplet.Tuplet      `yaml:"tuplet,omitempty"`
	TimeLine *time_line.TimeLine `yaml:"timeLine,omitempty"`
}

func (s *Symbol) IsNote() bool {
	return s.Note != nil
}

func (s *Symbol) IsEmbellishment() bool {
	return s.Note != nil && s.Note.Embellishment != nil
}

func (s *Symbol) IsValidNote() bool {
	if s.Note != nil {
		return s.Note.IsValid()
	}

	return false
}
