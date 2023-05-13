package symbols

import (
	"banduslib/internal/common"
)

//go:generate go run github.com/dmarkham/enumer -json -yaml -type=Accidental
//go:generate go run github.com/dmarkham/enumer -json -yaml -type=Tie

type Accidental uint

const (
	NoAccidental Accidental = iota
	Sharp
	Flat
	Natural
)

type Tie uint

const (
	NoTie Tie = iota
	Start
	End
)

type Note struct {
	Pitch         common.Pitch   `yaml:"pitch,omitempty"`
	Length        common.Length  `yaml:"length,omitempty"`
	Dots          uint8          `yaml:"dots,omitempty"`
	Accidental    Accidental     `yaml:"accidental,omitempty"`
	Fermata       bool           `yaml:"fermata,omitempty"`
	Tie           Tie            `yaml:"tie,omitempty"`
	Embellishment *Embellishment `yaml:"embellishment,omitempty"`
}

func (n *Note) IsValid() bool {
	if n.Pitch == common.NoPitch && n.Length == common.NoLength {
		return false
	}
	return true
}

// IsIncomplete returns true, if the note has no pitch and length but already other
// properties like an embellishment or an accidental
// this is the case, when the bww symbols which modify the note are preceding the note in
// bww code
func (n *Note) IsIncomplete() bool {
	if !n.HasPitchAndLength() {
		if n.Embellishment != nil {
			return true
		}
		if n.Accidental != NoAccidental {
			return true
		}
		if n.Tie != NoTie {
			return true
		}
	}
	return false
}

func (n *Note) HasPitchAndLength() bool {
	return n.Pitch != common.NoPitch && n.Length != common.NoLength
}

func (n *Note) IsOnlyAccidental() bool {
	if n.Pitch == common.NoPitch && n.Length == common.NoLength {
		if n.Accidental != NoAccidental {
			return true
		}
	}
	return false
}
