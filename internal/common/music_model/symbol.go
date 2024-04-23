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

func (s *Symbol) IsTimeline() bool {
	return s.TimeLine != nil
}

func (s *Symbol) IsTimelineStart() bool {
	return s.isTimelineBoundary(time_line.Start)
}

func (s *Symbol) IsTimelineEnd() bool {
	return s.isTimelineBoundary(time_line.End)
}

func (s *Symbol) isTimelineBoundary(boundary time_line.TimeLineBoundary) bool {
	if s.TimeLine != nil &&
		s.TimeLine.BoundaryType == boundary {
		return true
	}

	return false
}
