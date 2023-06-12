package symbols

import (
	"banduslib/internal/common"
	"banduslib/internal/common/music_model/symbols/accidental"
	"banduslib/internal/common/music_model/symbols/embellishment"
	"banduslib/internal/common/music_model/symbols/movement"
	"banduslib/internal/common/music_model/symbols/tie"
	"banduslib/internal/common/music_model/symbols/tuplet"
)

type Note struct {
	Pitch         common.Pitch                 `yaml:"pitch,omitempty"`
	Length        common.Length                `yaml:"length,omitempty"`
	Dots          uint8                        `yaml:"dots,omitempty"`
	Accidental    accidental.Accidental        `yaml:"accidental,omitempty"`
	Fermata       bool                         `yaml:"fermata,omitempty"`
	Tie           tie.Tie                      `yaml:"tie,omitempty"`
	Tuplet        *tuplet.Tuplet               `yaml:"tuplet,omitempty"`
	Embellishment *embellishment.Embellishment `yaml:"embellishment,omitempty"`
	// ExpandedEmbellishment is a list of pitches that makes up the embellishment
	ExpandedEmbellishment []common.Pitch     `yaml:"-"`
	Movement              *movement.Movement `yaml:"movement,omitempty"` // Piobairached movements
	Comment               string             `yaml:"comment,omitempty"`
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
		if n.Movement != nil {
			return true
		}
		if n.Accidental != accidental.NoAccidental {
			return true
		}
		if n.Tie != tie.NoTie {
			return true
		}
		if n.Tuplet != nil {
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
		if n.Accidental != accidental.NoAccidental {
			return true
		}
	}
	return false
}
