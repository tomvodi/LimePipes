package model

import "encoding/xml"

type Measure struct {
	XMLName xml.Name `xml:"measure"`
	Number  int      `xml:"number,attr"`
}
