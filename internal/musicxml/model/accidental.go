package model

import (
	"banduslib/internal/common/music_model/symbols/accidental"
	"encoding/xml"
	"strings"
)

type Accidental struct {
	XMLName xml.Name `xml:"accidental"`
	Value   string   `xml:",chardata"`
}

func NewAccidentalFromMusicModel(acc accidental.Accidental) *Accidental {
	if acc == accidental.NoAccidental {
		return nil
	}

	return &Accidental{
		XMLName: xml.Name{
			Local: "accidental",
		},
		Value: strings.ToLower(acc.String()),
	}
}
