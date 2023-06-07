package model

import (
	"banduslib/internal/common"
	"banduslib/internal/common/music_model/symbols"
	"encoding/xml"
	"github.com/rs/zerolog/log"
)

type Note struct {
	XMLName  xml.Name `xml:"note"`
	Grace    *Grace   `xml:"grace,omitempty"`
	Pitch    Pitch    `xml:"pitch"`
	Duration uint8    `xml:"duration,omitempty"`
	Voice    uint8    `xml:"voice,omitempty"`
	Type     string   `xml:"type"`
	Stem     string   `xml:"stem,omitempty"`
	Beams    []Beam   `xml:"beam,omitempty"`
}

func NotesFromMusicModel(note *symbols.Note, divisions uint8) []Note {
	var notes []Note

	if note.Embellishment != nil && note.ExpandedEmbellishment != nil {

		for i, pitch := range note.ExpandedEmbellishment {
			grace := Note{
				XMLName: xml.Name{
					Local: "note",
				},
				Grace: NewGrace(),
				Pitch: PitchFromMusicModel(pitch),
				Voice: 1,
				Type:  typeFromLength(common.Thirtysecond),
				Stem:  "up",
			}
			if len(note.ExpandedEmbellishment) > 1 {
				grace.Beams = embellishmentBeamsForPosition(i, len(note.ExpandedEmbellishment))
			}
			notes = append(notes, grace)
		}
	}

	xmlNote := Note{
		XMLName: xml.Name{
			Local: "note",
		},
		Pitch:    PitchFromMusicModel(note.Pitch),
		Duration: durationFromLength(note.Length, divisions),
		Voice:    1,
		Type:     typeFromLength(note.Length),
		Stem:     stemFromLength(note.Length),
	}
	notes = append(notes, xmlNote)

	return notes
}

func embellishmentBeamsForPosition(idx int, len int) []Beam {
	var bType BeamType
	if idx == 0 {
		bType = Begin
	} else if idx == len-1 {
		bType = End
	} else {
		bType = Continue
	}
	const beamCnt = 3
	beams := make([]Beam, beamCnt)
	for i := uint8(0); i < beamCnt; i++ {
		beams[i] = NewBeam(i+1, bType)
	}

	return beams
}

func typeFromLength(length common.Length) string {
	switch length {
	case common.Whole:
		return "whole"
	case common.Half:
		return "half"
	case common.Quarter:
		return "quarter"
	case common.Eighth:
		return "eighth"
	case common.Sixteenth:
		return "16th"
	case common.Thirtysecond:
		return "32nd"
	}

	return ""
}

func stemFromLength(length common.Length) string {
	if length == common.Whole {
		return ""
	}

	return "down"
}

func durationFromLength(length common.Length, divisions uint8) uint8 {
	maxDivisions := (255 / 4)
	if divisions > uint8(maxDivisions) {
		log.Error().Msgf("divisions can't be greater than %d", maxDivisions)
		return 255
	}

	switch length {
	case common.Whole:
		return 4 * divisions
	case common.Half:
		return 2 * divisions
	case common.Quarter:
		return 1 * divisions
	case common.Eighth:
		return divisions / 2
	case common.Sixteenth:
		return divisions / 4
	case common.Thirtysecond:
		return divisions / 8
	}

	log.Error().Msgf("length %s not supported for calculation of note duration", length.String())

	return divisions
}
