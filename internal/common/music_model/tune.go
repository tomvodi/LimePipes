package music_model

type Tune struct {
	Title      string     `yaml:"title"`
	Type       string     `yaml:"type,omitempty"`
	Composer   string     `yaml:"composer,omitempty"`
	Footer     []string   `yaml:"footer,omitempty"`
	Comments   []string   `yaml:"comments,omitempty"`
	InLineText []string   `yaml:"inLineText,omitempty"`
	Tempo      uint64     `yaml:"tempo"`
	Measures   []*Measure `yaml:"measures"`
}
