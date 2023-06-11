package model

import (
	"banduslib/internal/musicxml/model/fermata"
	"banduslib/internal/musicxml/model/tied"
	"encoding/xml"
)

type Notations struct {
	XMLName xml.Name         `xml:"notations"`
	Fermata *fermata.Fermata `xml:"fermata,omitempty"`
	Tied    *tied.Tied       `xml:"tied,omitempty"`
}

func NewNotations() *Notations {
	return &Notations{
		XMLName: xml.Name{
			Local: "notations",
		},
	}
}
