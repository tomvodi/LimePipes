package model

import "encoding/xml"

type Dot struct {
	XMLName xml.Name `xml:"dot"`
}

func NewDot() Dot {
	return Dot{
		XMLName: xml.Name{
			Local: "dot",
		},
	}
}
