package symbols

//go:generate go run github.com/dmarkham/enumer -json -yaml -type=Pitch
//go:generate go run github.com/dmarkham/enumer -json -yaml -type=Length

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

type Note struct {
	Pitch         Pitch          `yaml:"pitch,omitempty"`
	Length        Length         `yaml:"length,omitempty"`
	Dots          uint8          `yaml:"dots,omitempty"`
	Embellishment *Embellishment `yaml:"embellishment,omitempty"`
}

func (n *Note) IsValid() bool {
	if n.Pitch == NoPitch && n.Length == NoLength {
		return false
	}
	return true
}

func (n *Note) IsEmbellishmentOnly() bool {
	if n.Pitch == NoPitch && n.Length == NoLength {
		if n.Embellishment != nil {
			return true
		}
	}
	return false
}
