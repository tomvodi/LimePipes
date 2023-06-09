package model

import "encoding/xml"

type Rest struct {
	XMLName xml.Name `xml:"rest"`
}

func NewRest() *Rest {
	return &Rest{
		XMLName: xml.Name{
			Local: "rest",
		},
	}
}
