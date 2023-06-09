package model

import (
	"banduslib/internal/common"
	"banduslib/internal/common/music_model/symbols/accidental"
	"encoding/xml"
)

type Pitch struct {
	XMLName xml.Name `xml:"pitch"`
	Step    string   `xml:"step"`
	Alter   int8     `xml:"alter,omitempty"` // -1 flat, 1 sharp
	Octave  uint8    `xml:"octave"`
}

func PitchFromMusicModel(pitch common.Pitch, acc accidental.Accidental) *Pitch {
	switch pitch {
	case common.LowG:
		return createPitch("G", 4, 0, acc)
	case common.LowA:
		return createPitch("A", 4, 0, acc)
	case common.B:
		return createPitch("B", 4, 0, acc)
	case common.C:
		return createPitch("C", 5, 1, acc)
	case common.D:
		return createPitch("D", 5, 0, acc)
	case common.E:
		return createPitch("E", 5, 0, acc)
	case common.F:
		return createPitch("F", 5, 1, acc)
	case common.HighG:
		return createPitch("G", 5, 0, acc)
	case common.HighA:
		return createPitch("A", 5, 0, acc)
	}

	return &Pitch{}
}

func createPitch(
	step string,
	octave uint8,
	alter int8,
	acc accidental.Accidental,
) *Pitch {
	retPitch := &Pitch{
		XMLName: xml.Name{Local: "pitch"},
		Step:    step,
		Octave:  octave,
		Alter:   alter,
	}
	if acc == accidental.Sharp {
		retPitch.Alter += 1
	}
	if acc == accidental.Flat {
		retPitch.Alter -= 1
	}

	return retPitch
}
