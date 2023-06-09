package model

import "encoding/xml"

type Attributes struct {
	XMLName   xml.Name `xml:"attributes"`
	Divisions uint8    `xml:"divisions"`
	Key       Key      `xml:"key"`
	Time      *Time    `xml:"time,omitempty"`
}

func NewAttributesMinimal() *Attributes {
	return &Attributes{
		XMLName: xml.Name{
			Local: "attributes",
		},
	}
}

func NewAttributesWithKey(divisions uint8) *Attributes {
	return &Attributes{
		XMLName: xml.Name{
			Local: "attributes",
		},
		Divisions: divisions,
		Key:       Key{Fifths: 2},
	}
}
