package symbols

import (
	"banduslib/internal/common"
)

//go:generate go run github.com/dmarkham/enumer -json -yaml -type=EmbellishmentType
//go:generate go run github.com/dmarkham/enumer -json -yaml -type=EmbellishmentVariant
//go:generate go run github.com/dmarkham/enumer -json -yaml -type=EmbellishmentWeight

type EmbellishmentType uint

const (
	NoEmbellishment EmbellishmentType = iota
	SingleGrace
	Doubling
	Strike
	Grip
	Taorluath
	Bubbly
	GraceBirl
	ABirl
	Birl
	ThrowD
	Pele
	DoubleStrike
	TripleStrike
	GTripleStrike
	ThumbTripleStrike
	HalfTripleStrike
	DoubleGrace
)

type EmbellishmentVariant uint

const (
	NoVariant EmbellishmentVariant = iota
	G
	Half
	Thumb
)

type EmbellishmentWeight uint

const (
	NoWeight EmbellishmentWeight = iota
	Light
	Heavy
)

type Embellishment struct {
	Type EmbellishmentType `yaml:"type"`
	// Pitch is set for embellishments that have a pitch regardless of the melody note
	// following it (e.g. single g-grace)
	// Other embellishments have their pitch defined by the melody note following it
	// (e.g. doubling) because a d-doubling can only precede a d-melody note.
	// In these cases, Pitch is set to NoPitch
	Pitch common.Pitch `yaml:"pitch,omitempty"`

	Variant EmbellishmentVariant `yaml:"variant,omitempty"`

	Weight EmbellishmentWeight `yaml:"weight,omitempty"`
}
