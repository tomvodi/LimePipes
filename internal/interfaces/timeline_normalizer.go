package interfaces

import "banduslib/internal/common/music_model"

//go:generate mockgen -source timeline_normalizer.go -destination ./mocks/mock_timeline_normalizer.go.go

// TimelineNormalizer normalizes a bww timeline which can have for example a 2 of 4 repeat
// time line where a second repeat in part 2 will also be played in part 4 but won't show up
// in part 4. As this can not represented by musicxml, the TimelineNormalizer takes the
// symbols of the second repeat in part 2 and copies them to part 4.
type TimelineNormalizer interface {
	NormalizeTimeline(tune *music_model.Tune) (*music_model.Tune, error)
}
