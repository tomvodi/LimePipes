package tie

//go:generate go run github.com/dmarkham/enumer -json -yaml -type=Tie


type Tie uint

const (
	NoTie Tie = iota
	Start
	End
)
