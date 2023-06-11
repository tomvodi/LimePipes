package fermata

import "encoding/xml"

//go:generate go run github.com/dmarkham/enumer -json -yaml -transform=lower -type=Type

type Type uint8

const (
	Upright Type = iota
	Inverted
)

type Fermata struct {
	XMLName xml.Name `xml:"fermata"`
	Type    string   `xml:"type,attr"`
}

func NewFermata(fType Type) *Fermata {
	return &Fermata{
		XMLName: xml.Name{
			Local: "fermata",
		},
		Type: fType.String(),
	}
}
