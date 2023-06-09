package barline

import "encoding/xml"

//go:generate go run github.com/dmarkham/enumer -json -yaml -transform=lower -type=Direction

type Direction uint8

const (
	Forward Direction = iota
	Backward
)

type Repeat struct {
	XMLName   xml.Name `xml:"repeat"`
	Direction string   `xml:"direction,attr"`
}

func NewRepeat(dir Direction) *Repeat {
	return &Repeat{
		XMLName:   xml.Name{Local: "repeat"},
		Direction: dir.String(),
	}
}
