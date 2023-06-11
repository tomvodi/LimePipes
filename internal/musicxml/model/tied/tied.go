package tied

import "encoding/xml"

//go:generate go run github.com/dmarkham/enumer -json -yaml -transform=lower -type=Type

type Type uint8

const (
	Start Type = iota
	Stop
)

type Tied struct {
	XMLName xml.Name `xml:"tied"`
	Type    string   `xml:"type,attr"`
}

func NewTied(tType Type) *Tied {
	return &Tied{XMLName: xml.Name{
		Local: "tied",
	},
		Type: tType.String(),
	}
}
