package time_line

//go:generate go run github.com/dmarkham/enumer -json -yaml -type=TimeLineBoundary
//go:generate go run github.com/dmarkham/enumer -json -yaml -type=TimeLineType

type TimeLineBoundary uint8

const (
	NoBoundary TimeLineBoundary = iota
	Start
	End
)

type TimeLineType uint8

const (
	NoType TimeLineType = iota
	First
	Singling
	Second
	Doubling
	SecondOf2
	SecondOf3
	SecondOf4
	SecondOf2And4
	SecondOf5
	SecondOf6
	SecondOf7
	SecondOf8
	Bis
	Intro
)

type TimeLine struct {
	BoundaryType TimeLineBoundary `yaml:"boundaryType"`
	Type         TimeLineType     `yaml:"type,omitempty"`
}
