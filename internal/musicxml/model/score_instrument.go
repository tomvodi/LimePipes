package model

import "encoding/xml"

type ScoreInstrument struct {
	XMLName xml.Name `xml:"score-instrument"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"instrument-name"`
}
