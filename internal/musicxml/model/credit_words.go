package model

import "encoding/xml"

type CreditWords struct {
	XMLName xml.Name `xml:"credit-words"`
	Value   string   `xml:",chardata"`
}
