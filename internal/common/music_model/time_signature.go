package music_model

type TimeSignature struct {
	Beats    uint8 `yaml:"beats"`
	BeatType uint8 `yaml:"beatType"`
}
