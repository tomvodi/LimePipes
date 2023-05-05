package music_model

type Tune struct {
	Title    string     `yaml:"title"`
	Type     string     `yaml:"type,omitempty"`
	Composer string     `yaml:"composer,omitempty"`
	Footer   string     `yaml:"footer,omitempty"`
	Tempo    uint64     `yaml:"tempo"`
	Measures []*Measure `yaml:"measures"`
}
