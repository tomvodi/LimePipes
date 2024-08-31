package importtype

//go:generate go run github.com/dmarkham/enumer -type=Type -transform=snake

type Type int

const (
	BWW Type = iota
)
