package timeline_normalizer

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/interfaces"
)

type timelineNormalize struct {
}

func (t *timelineNormalize) NormalizeTimeline(tune *music_model.Tune) (*music_model.Tune, error) {
	return tune, nil
}

func NewTimelineNormalizer() interfaces.TimelineNormalizer {
	return &timelineNormalize{}
}
