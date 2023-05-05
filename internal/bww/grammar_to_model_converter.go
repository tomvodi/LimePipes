package bww

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/symbols"
	"fmt"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

func convertGrammarToModel(grammar *BwwDocument) ([]*music_model.Tune, error) {
	var tunes []*music_model.Tune

	for _, tune := range grammar.Tunes {
		newTune := &music_model.Tune{}
		if err := fillTuneWithParmeter(newTune, tune.Header.TuneParameter); err != nil {
			return nil, err
		}

		// TODO when tempo only of first tune, set to other tunes as well?
		if err := fillTunePartsFromStaves(newTune, tune.Body.Staffs); err != nil {
			return nil, err
		}
		tunes = append(tunes, newTune)
	}

	return tunes, nil
}

func fillTuneWithParmeter(tune *music_model.Tune, params []*TuneParameter) error {
	for _, param := range params {
		if param.Tempo != nil {
			tempo, err := strconv.ParseUint(param.Tempo.Tempo, 10, 64)
			if err != nil {
				return fmt.Errorf("failed parsing tune tempo: %s", err.Error())
			}

			tune.Tempo = tempo
		}

		if param.Description != nil {
			firstParam := param.Description.ParamList[0]
			text := param.Description.Text
			if firstParam == "T" {
				tune.Title = text
			}
			if firstParam == "Y" {
				tune.Type = text
			}
			if firstParam == "M" {
				tune.Composer = text
			}
			if firstParam == "F" {
				tune.Footer = text
			}
		}
	}

	return nil
}

func fillTunePartsFromStaves(tune *music_model.Tune, staves []*Staff) error {
	var measures []*music_model.Measure
	for _, stave := range staves {
		staveMeasures, err := getMeasuresFromStave(stave)
		if err != nil {
			return err
		}

		measures = append(measures, staveMeasures...)
	}
	tune.Measures = measures

	return nil
}

func getMeasuresFromStave(stave *Staff) ([]*music_model.Measure, error) {
	var measures []*music_model.Measure
	currMeasure := &music_model.Measure{}
	for _, symbol := range stave.Symbols {
		// if symbol bar or part start => new measure
		// currMeasure to return measures
		if symbol.Barline != nil ||
			symbol.PartStart != nil {
			measures = append(measures, currMeasure)
			currMeasure = &music_model.Measure{}

			continue
		}

		var lastSym *music_model.Symbol
		measSymLen := len(currMeasure.Symbols)
		if len(currMeasure.Symbols) > 0 {
			lastSym = currMeasure.Symbols[measSymLen-1]
		}
		newSym, err := appendStaffSymbolToMeasureSymbols(symbol, lastSym)
		if err != nil {
			return nil, err
		}
		if newSym != nil {
			currMeasure.Symbols = append(currMeasure.Symbols, newSym)
		}
	}

	measures = append(measures, currMeasure)
	return measures, nil
}

func appendStaffSymbolToMeasureSymbols(
	staffSym *StaffSymbols,
	lastSym *music_model.Symbol,
) (*music_model.Symbol, error) {
	newSym := &music_model.Symbol{}

	if staffSym.WholeNote != nil || staffSym.HalfNote != nil ||
		staffSym.QuarterNote != nil || staffSym.EighthNote != nil ||
		staffSym.SixteenthNote != nil || staffSym.ThirtysecondNote != nil {
		// add melody note to last note if it is an embellishment
		if lastSym != nil && lastSym.Note != nil && lastSym.Note.IsEmbellishmentOnly() {
			handleNote(staffSym, lastSym.Note)
			return nil, nil
		} else {
			newSym.Note = &symbols.Note{}
			handleNote(staffSym, newSym.Note)
			return newSym, nil
		}
	}
	if staffSym.SingleGrace != nil {
		newSym.Note = &symbols.Note{
			Embellishment: embellishmentForSingleGrace(staffSym.SingleGrace),
		}
		return newSym, nil
	}
	if staffSym.SingleDots != nil || staffSym.DoubleDots != nil {
		handleDots(staffSym, lastSym)
	}

	return nil, nil //fmt.Errorf("staff symbol %v not handled", staffSym)
}

func handleDots(staffSym *StaffSymbols, lastSym *music_model.Symbol) {
	var dotCount = uint8(0)
	var dotSym *string
	if staffSym.SingleDots != nil {
		dotCount = 1
		dotSym = staffSym.SingleDots
	}
	if staffSym.DoubleDots != nil {
		dotCount = 2
		dotSym = staffSym.DoubleDots
	}
	if dotCount > 0 {
		if lastSym != nil && lastSym.Note != nil && lastSym.Note.IsValid() {
			lastSym.Note.Dots = dotCount
		} else {
			log.Error().Msgf("dot symbol %s is not preceded by melody note", *dotSym)
		}
	}
}

func handleNote(staffSym *StaffSymbols, note *symbols.Note) {
	var token *string
	if staffSym.WholeNote != nil {
		token = staffSym.WholeNote
	}
	if staffSym.HalfNote != nil {
		token = staffSym.HalfNote
	}
	if staffSym.QuarterNote != nil {
		token = staffSym.QuarterNote
	}
	if staffSym.EighthNote != nil {
		token = staffSym.EighthNote
	}
	if staffSym.SixteenthNote != nil {
		token = staffSym.SixteenthNote
	}
	if staffSym.ThirtysecondNote != nil {
		token = staffSym.ThirtysecondNote
	}
	note.Length = lengthFromStaffNoteSuffix(token)
	note.Pitch = pitchFromStaffNotePrefix(token)
}

func embellishmentForSingleGrace(grace *string) *symbols.Embellishment {
	emb := &symbols.Embellishment{
		Type: symbols.SingleGrace,
	}

	if *grace == "ag" {
		emb.Pitch = symbols.LowA
	}
	if *grace == "bg" {
		emb.Pitch = symbols.B
	}
	if *grace == "cg" {
		emb.Pitch = symbols.C
	}
	if *grace == "dg" {
		emb.Pitch = symbols.D
	}
	if *grace == "eg" {
		emb.Pitch = symbols.E
	}
	if *grace == "fg" {
		emb.Pitch = symbols.F
	}
	if *grace == "gg" {
		emb.Pitch = symbols.HighG
	}
	if *grace == "tg" {
		emb.Pitch = symbols.HighA
	}
	return emb
}

func pitchFromStaffNotePrefix(note *string) symbols.Pitch {
	if strings.HasPrefix(*note, "LG") {
		return symbols.LowG
	}
	if strings.HasPrefix(*note, "LA") {
		return symbols.LowA
	}
	if strings.HasPrefix(*note, "B") {
		return symbols.B
	}
	if strings.HasPrefix(*note, "C") {
		return symbols.C
	}
	if strings.HasPrefix(*note, "D") {
		return symbols.D
	}
	if strings.HasPrefix(*note, "E") {
		return symbols.E
	}
	if strings.HasPrefix(*note, "F") {
		return symbols.F
	}
	if strings.HasPrefix(*note, "HG") {
		return symbols.HighG
	}
	if strings.HasPrefix(*note, "HA") {
		return symbols.HighA
	}

	return symbols.NoPitch
}
func lengthFromStaffNoteSuffix(note *string) symbols.Length {
	if strings.HasSuffix(*note, "16") {
		return symbols.Sixteenth
	}
	if strings.HasSuffix(*note, "32") {
		return symbols.Thirtysecond
	}
	if strings.HasSuffix(*note, "1") {
		return symbols.Whole
	}
	if strings.HasSuffix(*note, "2") {
		return symbols.Half
	}
	if strings.HasSuffix(*note, "4") {
		return symbols.Quarter
	}
	if strings.HasSuffix(*note, "8") {
		return symbols.Eigth
	}

	return symbols.NoLength
}
