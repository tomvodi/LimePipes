package music_model

import "fmt"

type TimeSignature struct {
	Beats    uint8 `yaml:"beats"`
	BeatType uint8 `yaml:"beatType"`
}

func (t *TimeSignature) String() string {
	return fmt.Sprintf("%d/%d", t.Beats, t.BeatType)
}
