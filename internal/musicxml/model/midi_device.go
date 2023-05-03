package model

import "encoding/xml"

type MidiDevice struct {
	XMLName xml.Name `xml:"midi-device"`
	Id      string   `xml:"id,attr"`
	Port    uint     `xml:"port,attr"`
}
