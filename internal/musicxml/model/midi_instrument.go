package model

import "encoding/xml"

type MidiInstrument struct {
	XMLName  xml.Name `xml:"midi-instrument"`
	Id       string   `xml:"id,attr"`
	Channel  uint     `xml:"midi-channel"`
	Programm uint     `xml:"midi-programm"`
	Volume   float32  `xml:"volume"`
	Pan      uint     `xml:"pan"`
}
