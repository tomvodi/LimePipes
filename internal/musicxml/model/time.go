package model

import (
	"banduslib/internal/common/music_model"
	"encoding/xml"
)

type Time struct {
	XMLName  xml.Name `xml:"time"`
	Beats    uint8    `xml:"beats"`
	BeatType uint8    `xml:"beat-type"`
}

func NewTime(time *music_model.TimeSignature) *Time {
	return &Time{
		XMLName: xml.Name{
			Local: "time",
		},
		Beats:    time.Beats,
		BeatType: time.BeatType,
	}
}
