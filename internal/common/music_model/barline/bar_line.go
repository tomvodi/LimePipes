package barline

//go:generate go run github.com/dmarkham/enumer -json -yaml -type=BarlineType
//go:generate go run github.com/dmarkham/enumer -json -yaml -type=TimelineType

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
	Segno
	Dalsegno
	Fine
	DacapoAlFine
)

type Barline struct {
	Type     BarlineType  `yaml:"type"`
	Timeline TimelineType `yaml:"timeline,omitempty"`
}
