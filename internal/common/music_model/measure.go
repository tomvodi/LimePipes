package music_model

type Measure struct {
	Time       *TimeSignature `yaml:"time,omitempty"`
	Symbols    []*Symbol      `yaml:"symbols,omitempty"`
	Comments   []string       `yaml:"comments,omitempty"`
	InLineText []string       `yaml:"inLineText,omitempty"`
}
