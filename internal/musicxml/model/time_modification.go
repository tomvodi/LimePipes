package model

import (
	"banduslib/internal/common/music_model/symbols/tuplet"
	"encoding/xml"
)

type TimeModification struct {
	XMLName     xml.Name `xml:"time-modification"`
	ActualNotes uint8    `xml:"actual-notes"`
	NormalNotes uint8    `xml:"normal-notes"`
}

func NewTimeModification(tpl *tuplet.Tuplet) *TimeModification {
	return &TimeModification{
		XMLName: xml.Name{
			Local: "time-modification",
		},
		ActualNotes: tpl.VisibleNotes,
		NormalNotes: tpl.PlayedNotes,
	}
}
