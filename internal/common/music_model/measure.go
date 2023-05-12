package music_model

type Measure struct {
	Time    *TimeSignature `yaml:"time,omitempty"`
	Symbols []*Symbol      `yaml:"symbols,omitempty"`
}
