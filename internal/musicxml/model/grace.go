package model

import "encoding/xml"

type Grace struct {
	XMLName xml.Name `xml:"grace"`
}

func NewGrace() *Grace {
	return &Grace{
		XMLName: xml.Name{
			Local: "grace",
		},
	}
}
