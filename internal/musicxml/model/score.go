package model

import "encoding/xml"

type Score struct {
	XMLName  xml.Name      `xml:"score-partwise"`
	Version  string        `xml:"version,attr"`
	PartList ScorePartList `xml:"part-list"`
	Part     Part          `xml:"part"`
	Credits  []Credit      `xml:"credit"`
}
