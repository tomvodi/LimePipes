package model

import "encoding/xml"

type Score struct {
	XMLName  xml.Name `xml:"score-partwise"`
	PartList ScorePartList
	Credits  []Credit
}
