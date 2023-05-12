package symbols

//go:generate go run github.com/dmarkham/enumer -json -yaml -type=EmbellishmentType

type EmbellishmentType uint

const (
	NoEmbellishment EmbellishmentType = iota
	SingleGrace
	Doubling
	HalfDoubling
	ThumbDoubling
	Strike
	GStrike
	ThumbStrike
	HalfStrike
	Grip
	HalfGrip
	ThumbGrip
	GGrip
	Taorluath
	Bubbly
	GraceBirl
	ABirl
	Birl
	ThrowD
	HeavyThrowD
	Pele
	ThumbPele
	HalfPele
	DoubleStrike
	GDoubleStrike
	ThumbDoubleStrike
	HalfDoubleStrike
	TripleStrike
	GTripleStrike
	ThumbTripleStrike
	HalfTripleStrike
	DDoubleGrace
	EDoubleGrace
	FDoubleGrace
	GDoubleGrace
	ThumbDoubleGrace
)

type Embellishment struct {
	Type EmbellishmentType `yaml:"type"`
	// Pitch is set for embellishments that have a pitch regardless of the melody note
	// following it (e.g. single g-grace)
	// Other embellishments have their pitch defined by the melody note following it
	// (e.g. doubling) because a d-doubling can only precede a d-melody note.
	// In these cases, Pitch is set to NoPitch
	Pitch Pitch `yaml:"pitch,omitempty"`
}
