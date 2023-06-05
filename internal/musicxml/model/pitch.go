package model

import (
	"banduslib/internal/common"
	"encoding/xml"
)

type Pitch struct {
	XMLName xml.Name `xml:"pitch"`
	Step    string   `xml:"step"`
	Alter   uint8    `xml:"alter,omitempty"` // -1 flat, 1 sharp
	Octave  uint8    `xml:"octave"`
}

func PitchFromMusicModel(pitch common.Pitch) Pitch {
	switch pitch {
	case common.LowG:
		return createPitch("G", 4, 0)
	case common.LowA:
		return createPitch("A", 4, 0)
	case common.B:
		return createPitch("B", 4, 0)
	case common.C:
		return createPitch("C", 5, 1)
	case common.D:
		return createPitch("D", 5, 0)
	case common.E:
		return createPitch("E", 5, 0)
	case common.F:
		return createPitch("F", 5, 1)
	case common.HighG:
		return createPitch("G", 5, 0)
	case common.HighA:
		return createPitch("A", 5, 0)
	}

	return Pitch{}
}

func createPitch(step string, octave uint8, alter uint8) Pitch {
	return Pitch{
		XMLName: xml.Name{Local: "pitch"},
		Step:    step,
		Octave:  octave,
		Alter:   alter,
	}
}
