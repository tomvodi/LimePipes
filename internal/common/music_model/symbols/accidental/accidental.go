package accidental

//go:generate go run github.com/dmarkham/enumer -json -yaml -type=Accidental


type Accidental uint

const (
	NoAccidental Accidental = iota
	Sharp
	Flat
	Natural
)
