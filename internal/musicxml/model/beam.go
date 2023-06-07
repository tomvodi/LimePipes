package model

import "encoding/xml"

//go:generate go run github.com/dmarkham/enumer -json -yaml -transform=lower -type=BeamType

type BeamType uint

const (
	Begin BeamType = iota
	Continue
	End
)

type Beam struct {
	XMLName xml.Name `xml:"beam"`
	Number  uint8    `xml:"number,attr"`
	Value   string   `xml:",chardata"`
}

func NewBeam(nr uint8, bType BeamType) Beam {
	return Beam{
		XMLName: xml.Name{
			Local: "beam",
		},
		Number: nr,
		Value:  bType.String(),
	}
}
