package common

//go:generate go run github.com/dmarkham/enumer -json -yaml -type=Pitch
//go:generate go run github.com/dmarkham/enumer -json -yaml -type=Length

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
