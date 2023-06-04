package model

import "encoding/xml"

type MidiInstrument struct {
	XMLName xml.Name `xml:"midi-instrument"`
	Id      string   `xml:"id,attr"`
	Channel uint     `xml:"midi-channel"`
	Program uint     `xml:"midi-program"`
	Volume  float32  `xml:"volume"`
	Pan     uint     `xml:"pan"`
}
