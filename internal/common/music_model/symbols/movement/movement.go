package movement

import (
	"banduslib/internal/common"
)

//go:generate go run github.com/dmarkham/enumer -json -yaml -type=Type
//go:generate go run github.com/dmarkham/enumer -json -yaml -type=Variant

type Type uint

const (
	NoMovement Type = iota
	Cadence
	Embari
	Endari
	Chedari
	Hedari
	Dili
	Tra
	Edre
	Dare
	CheCheRe
	Grip
	Deda
	Enbain
	Otro
	Odro
	Adeda
	EchoBeat
	Darodo
	Hiharin
	Rodin
	Chelalho
	Din
	Lemluath
	Taorluath
	Crunluath
)

type Variant uint

const (
	NoVariant Variant = iota
	Half
	G
	Thumb
	LongLowG
)

type Movement struct {
	Type    Type           `yaml:"type"`
	Pitches []common.Pitch `yaml:"pitches,omitempty"` // Only used by cadence
	Fermata bool           `yaml:"fermata,omitempty"` // Only used by cadence
	Variant Variant        `yaml:"variant,omitempty"`
	// Bww symbols for piobaireachd are written with a pitch e.g. edred which
	// indicates that a B note should precede this note.
	// Nevertheless, in some cases, a note with another pitch is preceding the movement
	// so the pitch hint will make sure that no information gets lost in importing bww
	// files
	PitchHint           common.Pitch `yaml:"pitchHint,omitempty"`
	AdditionalPitchHint common.Pitch `yaml:"additionalPitchHint,omitempty"`
	// Pitch, if the movement has a distinct pitch like echo notes
	Pitch common.Pitch `yaml:"pitch,omitempty"`
	// Abbreviate is true, when the movement should show as its abbreviation symbol and not
	// as every grace note
	Abbreviate bool `yaml:"abbreviate,omitempty"`

	Breabach bool `yaml:"breabach,omitempty"`
	AMach    bool `yaml:"AMach,omitempty"`
}
