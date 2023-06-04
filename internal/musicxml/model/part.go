package model

import "encoding/xml"

type Part struct {
	XMLName  xml.Name  `xml:"part"`
	Id       string    `xml:"id,attr"`
	Measures []Measure `xml:"measure"`
}
