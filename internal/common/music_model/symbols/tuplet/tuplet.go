package tuplet

//go:generate go run github.com/dmarkham/enumer -json -yaml -type=TupletBoundary
//go:generate go run github.com/dmarkham/enumer -json -yaml -type=TupletType

type TupletBoundary uint8

const (
	NoBoundary TupletBoundary = iota
	Start
	End
)

type TupletType uint8

const (
	NoType TupletType = iota
	Type23
	Type32
	Type43
	Type46
	Type53
	Type54
	Type64
	Type74
	Type76
)

type Tuplet struct {
	BoundaryType TupletBoundary `yaml:"boundaryType"`
	VisibleNotes uint8          `yaml:"visibleNotes"`
	PlayedNotes  uint8          `yaml:"playedNotes"`
}

func NewTuplet(bound TupletBoundary, ttype TupletType) *Tuplet {
	tp := &Tuplet{
		BoundaryType: bound,
	}
	tp.VisibleNotes, tp.PlayedNotes = notesConfigFromType(ttype)

	return tp
}

// notesConfigFromType returns the visible notes and played notes for a given type
func notesConfigFromType(ttype TupletType) (uint8, uint8) {
	switch ttype {
	case Type23:
		return 2, 3
	case Type32:
		return 3, 2
	case Type43:
		return 4, 3
	case Type46:
		return 4, 6
	case Type53:
		return 5, 3
	case Type54:
		return 5, 4
	case Type64:
		return 6, 4
	case Type74:
		return 7, 4
	case Type76:
		return 7, 6
	default:
		return 0, 0
	}
}
