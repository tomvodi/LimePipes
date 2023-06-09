package model

import (
	"banduslib/internal/musicxml/model/barline"
	"encoding/xml"
)

type Measure struct {
	XMLName    xml.Name          `xml:"measure"`
	Number     int               `xml:"number,attr"`
	Barlines   []barline.Barline `xml:"barline"`
	Attributes *Attributes       `xml:"attributes,omitempty"`
	Notes      []Note            `xml:"note"`
}
