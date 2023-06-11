package model

import (
	"banduslib/internal/musicxml/model/fermata"
	"encoding/xml"
)

type Notations struct {
	XMLName xml.Name         `xml:"notations"`
	Fermata *fermata.Fermata `xml:"fermata,omitempty"`
}

func NewNotations() *Notations {
	return &Notations{
		XMLName: xml.Name{
			Local: "notations",
		},
	}
}
