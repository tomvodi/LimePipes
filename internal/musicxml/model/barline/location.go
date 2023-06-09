package barline

//go:generate go run github.com/dmarkham/enumer -json -yaml -transform=lower -type=Location

type Location uint8

const (
	Left Location = iota
	Right
)
