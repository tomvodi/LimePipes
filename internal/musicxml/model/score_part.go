package model

import "encoding/xml"

type ScorePart struct {
	XMLName        xml.Name        `xml:"score-part"`
	Id             string          `xml:"id,attr"`
	Name           string          `xml:"part-name"`
	Instrument     ScoreInstrument `xml:"score-instrument"`
	MidiDevice     MidiDevice      `xml:"midi-device"`
	MidiInstrument MidiInstrument  `xml:"midi-instrument"`
}
