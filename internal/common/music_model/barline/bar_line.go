package barline

//go:generate go run github.com/dmarkham/enumer -json -yaml -type=BarlineType
//go:generate go run github.com/dmarkham/enumer -json -yaml -type=TimelineType
//go:generate go run github.com/dmarkham/enumer -json -yaml -type=SegnoType
//go:generate go run github.com/dmarkham/enumer -json -yaml -type=DacapoType

type BarlineType uint

const (
	Regular BarlineType = iota
	Heavy
	HeavyHeavy
	LightHeavy
	HeavyLight
)

type TimelineType uint

const (
	NoTimeline TimelineType = iota
	Repeat                  // a normal |: ... :| repeat lines
)

type SegnoType uint

const (
	NoSegnoType SegnoType = iota
	Segno
	Dalsegno
)

type DacapoType uint

const (
	NoDacapoType DacapoType = iota
	Fine
	DacapoAlFine
)

type Barline struct {
	Type       BarlineType  `yaml:"type"`
	Timeline   TimelineType `yaml:"timeline,omitempty"`
	SegnoType  SegnoType    `yaml:"segnoType,omitempty"`
	DacapoType DacapoType   `yaml:"dacapoType,omitempty"`
}

func (b *Barline) IsHeavy() bool {
	if b.Type == Heavy ||
		b.Type == HeavyLight ||
		b.Type == LightHeavy ||
		b.Type == HeavyHeavy {
		return true
	}

	return false
}
