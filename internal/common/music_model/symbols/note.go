package symbols

//go:generate go run github.com/dmarkham/enumer -json -yaml -type=Pitch
//go:generate go run github.com/dmarkham/enumer -json -yaml -type=Length
//go:generate go run github.com/dmarkham/enumer -json -yaml -type=Accidental

type Pitch uint

const (
	NoPitch Pitch = iota
	LowG
	LowA
	B
	C
	D
	E
	F
	HighG
	HighA
)

type Length uint

const (
	NoLength Length = iota
	Whole
	Half
	Quarter
	Eighth
	Sixteenth
	Thirtysecond
)

type Accidental uint

const (
	NoAccidental Accidental = iota
	Sharp
	Flat
	Natural
)

type Note struct {
	Pitch         Pitch          `yaml:"pitch,omitempty"`
	Length        Length         `yaml:"length,omitempty"`
	Dots          uint8          `yaml:"dots,omitempty"`
	Accidental    Accidental     `yaml:"accidental,omitempty"`
	Fermata       bool           `yaml:"fermata,omitempty"`
	Embellishment *Embellishment `yaml:"embellishment,omitempty"`
}

func (n *Note) IsValid() bool {
	if n.Pitch == NoPitch && n.Length == NoLength {
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
	}
	return false
}

func (n *Note) HasPitchAndLength() bool {
	return n.Pitch != NoPitch && n.Length != NoLength
}

func (n *Note) IsOnlyAccidental() bool {
	if n.Pitch == NoPitch && n.Length == NoLength {
		if n.Accidental != NoAccidental {
			return true
		}
	}
	return false
}
