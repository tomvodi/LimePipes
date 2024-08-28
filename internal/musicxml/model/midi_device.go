package model

import "encoding/xml"

type MidiDevice struct {
	XMLName xml.Name `xml:"midi-device"`
	ID      string   `xml:"id,attr"`
	Port    uint     `xml:"port,attr"`
}
