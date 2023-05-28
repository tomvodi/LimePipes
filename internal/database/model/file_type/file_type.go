package file_type

//go:generate go run github.com/dmarkham/enumer -sql -type=Type

type Type uint8

const (
	MusicModel Type = iota + 1
	MusicXml
	Bww
)
