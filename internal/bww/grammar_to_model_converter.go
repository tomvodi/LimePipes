package bww

import (
	"banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/symbols"
	"banduslib/internal/common/music_model/symbols/tuplet"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"strconv"
	"strings"
)

func convertGrammarToModel(grammar *BwwDocument) ([]*music_model.Tune, error) {
	var tunes []*music_model.Tune

	for _, tune := range grammar.Tunes {
		newTune := &music_model.Tune{}
		if err := fillTuneWithParameter(newTune, tune.Header.TuneParameter); err != nil {
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

func fillTuneWithParameter(tune *music_model.Tune, params []*TuneParameter) error {
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
	for _, staffSym := range stave.Symbols {
		// if staffSym bar or part start => new measure
		// currMeasure to return measures
		if staffSym.Barline != nil ||
			staffSym.PartStart != nil {
			measures = cleanupAndAppendMeasure(measures, currMeasure)
			currMeasure = &music_model.Measure{}

			continue
		}

		if staffSym.TimeSig != nil {
			setMeasureTimeSig(currMeasure, staffSym.TimeSig)
			continue
		}

		// triplets in old format appear only after last melody note,
		// so it is handled here
		if staffSym.Triplets != nil {
			err := handleTriplet(currMeasure, *staffSym.Triplets)
			if err != nil {
				return nil, err
			}
		}

		var lastSym *music_model.Symbol
		measSymLen := len(currMeasure.Symbols)
		if len(currMeasure.Symbols) > 0 {
			lastSym = currMeasure.Symbols[measSymLen-1]
		}
		newSym, err := appendStaffSymbolToMeasureSymbols(staffSym, lastSym, currMeasure)
		if err != nil {
			return nil, err
		}
		if newSym != nil {
			currMeasure.Symbols = append(currMeasure.Symbols, newSym)
		}
	}

	measures = cleanupAndAppendMeasure(measures, currMeasure)
	return measures, nil
}

func handleTriplet(measure *music_model.Measure, sym string) error {

	if len(measure.Symbols) == 0 {
		return fmt.Errorf("triplet symbol %s does not follow any note", sym)
	}
	if len(measure.Symbols) < 3 {
		return fmt.Errorf("triplet symbol %s must follow at least 3 notes", sym)
	}

	var last3SymsAreNotes = true
	lastIndex := len(measure.Symbols) - 1
	for i := lastIndex; i > lastIndex-3; i-- {
		currSym := measure.Symbols[i]
		if currSym.Note == nil {
			last3SymsAreNotes = false
			break
		}
		if !currSym.Note.IsValid() {
			last3SymsAreNotes = false
			break
		}
	}
	if !last3SymsAreNotes {
		return fmt.Errorf("triplet symbol %s must follow at least 3 notes", sym)
	}

	tripletStartIdx := lastIndex - 2
	hasSymbolsBeforeTriplet := tripletStartIdx > 0
	tupletHasAlreadyAStartSymbol := false
	if hasSymbolsBeforeTriplet {
		symBeforeTriplet := measure.Symbols[tripletStartIdx-1]
		if symBeforeTriplet.Tuplet != nil &&
			symBeforeTriplet.Tuplet.BoundaryType == tuplet.Start {
			tupletHasAlreadyAStartSymbol = true
		}
	}
	if !tupletHasAlreadyAStartSymbol {
		tripletStartSym := newIrregularGroup(tuplet.Start, tuplet.Type32)
		measure.Symbols = append(
			measure.Symbols[:tripletStartIdx+1],
			measure.Symbols[tripletStartIdx:]...,
		)
		measure.Symbols[tripletStartIdx] = tripletStartSym
	}

	measure.Symbols = append(measure.Symbols, newIrregularGroup(tuplet.End, tuplet.Type32))

	return nil
}

func cleanupAndAppendMeasure(
	measures []*music_model.Measure,
	measure *music_model.Measure,
) []*music_model.Measure {
	cleanupMeasure(measure)
	if len(measure.Symbols) == 0 {
		measure.Symbols = nil
	}
	return append(measures, measure)
}

// cleanupMeasure removes invalid symbols from the measure
// this may be the case for the accidentals at the beginning of the measure which are
// indicating the key of the measure. For bagpipes the key is always sharpf sharpc,
// so we delete these symbols here.
func cleanupMeasure(meas *music_model.Measure) {
	for _, symbol := range meas.Symbols {
		if symbol.Note == nil {
			continue
		}

		if symbol.Note.IsOnlyAccidental() {
			idx := symbolIndexOf(meas.Symbols, symbol)
			if idx == -1 {
				log.Error().Msgf("symbol index could not be found for cleanup of measure")
			} else {
				meas.Symbols = removeSymbol(meas.Symbols, idx)
			}
		}
	}
}

func setMeasureTimeSig(measure *music_model.Measure, timeSigSym *string) {
	timeSig := timeSigFromSymbol(timeSigSym)
	measure.Time = timeSig
}

func timeSigFromSymbol(sym *string) *music_model.TimeSignature {
	if *sym == "C" {
		return &music_model.TimeSignature{
			Beats:    4,
			BeatType: 4,
		}
	}

	if *sym == "C_" {
		return &music_model.TimeSignature{
			Beats:    2,
			BeatType: 2,
		}
	}

	parts := strings.Split(*sym, "_")
	if len(parts) != 2 {
		log.Error().Msgf("time signature symbol %s can't be parsed", *sym)
		return nil
	}

	beat, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		log.Error().Msgf("failed parsing time signature beats part %s: %s", parts[0], err.Error())
		return nil
	}

	beatTime, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		log.Error().Msgf("failed parsing time signature beat type part %s: %s", parts[1], err.Error())
		return nil
	}

	return &music_model.TimeSignature{
		Beats:    uint8(beat),
		BeatType: uint8(beatTime),
	}
}

func symbolIndexOf(symbols []*music_model.Symbol, findSym *music_model.Symbol) int {
	for i, symbol := range symbols {
		if reflect.DeepEqual(symbol, findSym) {
			return i
		}
	}
	return -1
}

func removeSymbol(symbols []*music_model.Symbol, idx int) []*music_model.Symbol {
	return append(symbols[:idx], symbols[idx+1:]...)
}

func appendStaffSymbolToMeasureSymbols(
	staffSym *StaffSymbols,
	lastSym *music_model.Symbol,
	currentMeasure *music_model.Measure,
) (*music_model.Symbol, error) {
	newSym := &music_model.Symbol{}

	if staffSym.WholeNote != nil || staffSym.HalfNote != nil ||
		staffSym.QuarterNote != nil || staffSym.EighthNote != nil ||
		staffSym.SixteenthNote != nil || staffSym.ThirtysecondNote != nil {
		// add melody note to last note if it is an embellishment
		if lastSym != nil && lastSym.Note != nil && lastSym.Note.IsIncomplete() {
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
	if staffSym.TieStart != nil {
		newSym.Note = &symbols.Note{
			Tie: symbols.Start,
		}
		return newSym, nil
	}
	if staffSym.TieEnd != nil {
		// TODO: check if tie start note has same pitch and if a single ^te could
		// be an old tie format for two E notes
		if lastSym != nil && lastSym.Note != nil {
			lastSym.Note.Tie = symbols.End
		}
	}
	if staffSym.TieOld != nil {
		return nil, fmt.Errorf("old tie format not supported %s", *staffSym.TieOld)
	}
	if staffSym.Flat != nil {
		return handleAccidential(symbols.Flat), nil
	}
	if staffSym.Natural != nil {
		return handleAccidential(symbols.Natural), nil
	}
	if staffSym.Sharp != nil {
		return handleAccidential(symbols.Sharp), nil
	}
	if staffSym.Doubling != nil {
		return handleEmbellishment(symbols.Doubling)
	}
	if staffSym.HalfDoubling != nil {
		return handleEmbellishmentVariant(symbols.Doubling, symbols.Half, symbols.NoWeight)
	}
	if staffSym.ThumbDoubling != nil {
		return handleEmbellishmentVariant(symbols.Doubling, symbols.Thumb, symbols.NoWeight)
	}
	if staffSym.Grip != nil {
		return handleEmbellishment(symbols.Grip)
	}
	if staffSym.GGrip != nil {
		return handleEmbellishmentVariant(symbols.Grip, symbols.G, symbols.NoWeight)
	}
	if staffSym.ThumbGrip != nil {
		return handleEmbellishmentVariant(symbols.Grip, symbols.Thumb, symbols.NoWeight)
	}
	if staffSym.HalfGrip != nil {
		return handleEmbellishmentVariant(symbols.Grip, symbols.Half, symbols.NoWeight)
	}
	if staffSym.Taorluath != nil {
		return handleEmbellishment(symbols.Taorluath)
	}
	if staffSym.Bubbly != nil {
		return handleEmbellishment(symbols.Bubbly)
	}
	if staffSym.ThrowD != nil {
		return handleEmbellishmentVariant(symbols.ThrowD, symbols.NoVariant, symbols.Light)
	}
	if staffSym.HeavyThrowD != nil {
		return handleEmbellishment(symbols.ThrowD)
	}
	if staffSym.Birl != nil {
		return handleEmbellishment(symbols.Birl)
	}
	if staffSym.ABirl != nil {
		return handleEmbellishment(symbols.ABirl)
	}
	if staffSym.Strike != nil {
		return handleEmbellishment(symbols.Strike)
	}
	if staffSym.GStrike != nil {
		return handleEmbellishmentVariant(symbols.Strike, symbols.G, symbols.NoWeight)
	}
	if staffSym.LightGStrike != nil {
		return handleEmbellishmentVariant(symbols.Strike, symbols.G, symbols.Light)
	}
	if staffSym.LightDoubleStrike != nil {
		return handleEmbellishmentVariant(symbols.DoubleStrike, symbols.NoVariant, symbols.Light)
	}
	if staffSym.DoubleStrike != nil {
		return handleEmbellishmentVariant(symbols.DoubleStrike, symbols.NoVariant, symbols.NoWeight)
	}
	if staffSym.LightGDoubleStrike != nil {
		return handleEmbellishmentVariant(symbols.DoubleStrike, symbols.G, symbols.Light)
	}
	if staffSym.GDoubleStrike != nil {
		return handleEmbellishmentVariant(symbols.DoubleStrike, symbols.G, symbols.NoWeight)
	}
	if staffSym.LightThumbDoubleStrike != nil {
		return handleEmbellishmentVariant(symbols.DoubleStrike, symbols.Thumb, symbols.Light)
	}
	if staffSym.ThumbDoubleStrike != nil {
		return handleEmbellishmentVariant(symbols.DoubleStrike, symbols.Thumb, symbols.NoWeight)
	}
	if staffSym.LightHalfDoubleStrike != nil {
		return handleEmbellishmentVariant(symbols.DoubleStrike, symbols.Half, symbols.Light)
	}
	if staffSym.HalfDoubleStrike != nil {
		return handleEmbellishmentVariant(symbols.DoubleStrike, symbols.Half, symbols.NoWeight)
	}
	if staffSym.LightTripleStrike != nil {
		return handleEmbellishmentVariant(symbols.TripleStrike, symbols.NoVariant, symbols.Light)
	}
	if staffSym.TripleStrike != nil {
		return handleEmbellishmentVariant(symbols.TripleStrike, symbols.NoVariant, symbols.NoWeight)
	}
	if staffSym.LightGTripleStrike != nil {
		return handleEmbellishmentVariant(symbols.TripleStrike, symbols.G, symbols.Light)
	}
	if staffSym.GTripleStrike != nil {
		return handleEmbellishmentVariant(symbols.TripleStrike, symbols.G, symbols.NoWeight)
	}
	if staffSym.LightThumbTripleStrike != nil {
		return handleEmbellishmentVariant(symbols.TripleStrike, symbols.Thumb, symbols.Light)
	}
	if staffSym.ThumbTripleStrike != nil {
		return handleEmbellishmentVariant(symbols.TripleStrike, symbols.Thumb, symbols.NoWeight)
	}
	if staffSym.LightHalfTripleStrike != nil {
		return handleEmbellishmentVariant(symbols.TripleStrike, symbols.Half, symbols.Light)
	}
	if staffSym.HalfTripleStrike != nil {
		return handleEmbellishmentVariant(symbols.TripleStrike, symbols.Half, symbols.NoWeight)
	}
	if staffSym.DDoubleGrace != nil {
		return handleDoubleGrace(common.D)
	}
	if staffSym.EDoubleGrace != nil {
		return handleDoubleGrace(common.E)
	}
	if staffSym.FDoubleGrace != nil {
		return handleDoubleGrace(common.F)
	}
	if staffSym.GDoubleGrace != nil {
		return handleDoubleGrace(common.HighG)
	}
	if staffSym.ThumbDoubleGrace != nil {
		return handleDoubleGrace(common.HighA)
	}
	if staffSym.HalfStrike != nil {
		return handleEmbellishmentVariant(symbols.Strike, symbols.Half, symbols.NoWeight)
	}
	if staffSym.LightHalfStrike != nil {
		return handleEmbellishmentVariant(symbols.Strike, symbols.Half, symbols.Light)
	}
	if staffSym.ThumbStrike != nil {
		return handleEmbellishmentVariant(symbols.Strike, symbols.Thumb, symbols.NoWeight)
	}
	if staffSym.LightThumbStrike != nil {
		return handleEmbellishmentVariant(symbols.Strike, symbols.Thumb, symbols.Light)
	}
	if staffSym.Pele != nil {
		return handleEmbellishment(symbols.Pele)
	}
	if staffSym.LightPele != nil {
		return handleEmbellishmentVariant(symbols.Pele, symbols.NoVariant, symbols.Light)
	}
	if staffSym.ThumbPele != nil {
		return handleEmbellishmentVariant(symbols.Pele, symbols.Thumb, symbols.NoWeight)
	}
	if staffSym.LightThumbPele != nil {
		return handleEmbellishmentVariant(symbols.Pele, symbols.Thumb, symbols.Light)
	}
	if staffSym.HalfPele != nil {
		return handleEmbellishmentVariant(symbols.Pele, symbols.Half, symbols.NoWeight)
	}
	if staffSym.IrregularGroupStart != nil {
		ttype := tupletTypeFromSymbol(staffSym.IrregularGroupStart)
		return handleIrregularGroup(tuplet.Start, ttype)
	}
	if staffSym.IrregularGroupEnd != nil {
		ttype := tupletTypeFromSymbol(staffSym.IrregularGroupEnd)
		// handling old style ties on E (^3e) but ignoring an error
		if ttype == tuplet.Type32 {
			_ = handleTriplet(currentMeasure, "^3e")
			return nil, nil
		} else {
			return handleIrregularGroup(tuplet.End, ttype)
		}
	}
	if staffSym.LightHalfPele != nil {
		return handleEmbellishmentVariant(symbols.Pele, symbols.Half, symbols.Light)
	}
	if staffSym.GBirl != nil ||
		staffSym.ThumbBirl != nil {
		return handleEmbellishment(symbols.GraceBirl)
	}
	if staffSym.Fermata != nil {
		if lastSym != nil && lastSym.Note != nil && lastSym.Note.HasPitchAndLength() {
			lastSym.Note.Fermata = true
		}
	}
	if staffSym.Rest != nil {
		newSym.Rest = &symbols.Rest{
			Length: lengthFromSuffix(staffSym.Rest),
		}
		return newSym, nil
	}

	return nil, nil // fmt.Errorf("staff symbol %v not handled", staffSym)
}

func handleEmbellishment(
	emb symbols.EmbellishmentType,
) (*music_model.Symbol, error) {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Embellishment: &symbols.Embellishment{
				Type: emb,
			},
		},
	}, nil
}

func handleDoubleGrace(pitch common.Pitch) (*music_model.Symbol, error) {
	doubleG, err := handleEmbellishment(symbols.DoubleGrace)
	doubleG.Note.Embellishment.Pitch = pitch
	return doubleG, err
}

func handleEmbellishmentVariant(
	emb symbols.EmbellishmentType,
	variant symbols.EmbellishmentVariant,
	weight symbols.EmbellishmentWeight,
) (*music_model.Symbol, error) {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Embellishment: &symbols.Embellishment{
				Type:    emb,
				Variant: variant,
				Weight:  weight,
			},
		},
	}, nil
}

func handleIrregularGroup(
	boundary tuplet.TupletBoundary,
	ttype tuplet.TupletType,
) (*music_model.Symbol, error) {
	return newIrregularGroup(boundary, ttype), nil
}

func newIrregularGroup(boundary tuplet.TupletBoundary,
	ttype tuplet.TupletType,
) *music_model.Symbol {
	tpl := tuplet.NewTuplet(boundary, ttype)
	return &music_model.Symbol{
		Tuplet: tpl,
	}
}

func tupletTypeFromSymbol(sym *string) tuplet.TupletType {
	if strings.HasPrefix(*sym, "^2") {
		return tuplet.Type23
	}
	if strings.HasPrefix(*sym, "^3") {
		return tuplet.Type32
	}
	if strings.HasPrefix(*sym, "^43") {
		return tuplet.Type43
	}
	if strings.HasPrefix(*sym, "^46") {
		return tuplet.Type46
	}
	if strings.HasPrefix(*sym, "^53") {
		return tuplet.Type53
	}
	if strings.HasPrefix(*sym, "^54") {
		return tuplet.Type54
	}
	if strings.HasPrefix(*sym, "^64") {
		return tuplet.Type64
	}
	if strings.HasPrefix(*sym, "^74") {
		return tuplet.Type74
	}
	if strings.HasPrefix(*sym, "^76") {
		return tuplet.Type76
	}

	log.Error().Msgf("tuplet symbold %s not handled", *sym)
	return tuplet.NoType
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
	note.Length = lengthFromSuffix(token)
	note.Pitch = pitchFromStaffNotePrefix(token)
}

func handleAccidential(acc symbols.Accidental) *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Accidental: acc,
		},
	}
}

func embellishmentForSingleGrace(grace *string) *symbols.Embellishment {
	emb := &symbols.Embellishment{
		Type: symbols.SingleGrace,
	}

	if *grace == "ag" {
		emb.Pitch = common.LowA
	}
	if *grace == "bg" {
		emb.Pitch = common.B
	}
	if *grace == "cg" {
		emb.Pitch = common.C
	}
	if *grace == "dg" {
		emb.Pitch = common.D
	}
	if *grace == "eg" {
		emb.Pitch = common.E
	}
	if *grace == "fg" {
		emb.Pitch = common.F
	}
	if *grace == "gg" {
		emb.Pitch = common.HighG
	}
	if *grace == "tg" {
		emb.Pitch = common.HighA
	}
	return emb
}

func pitchFromStaffNotePrefix(note *string) common.Pitch {
	if strings.HasPrefix(*note, "LG") {
		return common.LowG
	}
	if strings.HasPrefix(*note, "LA") {
		return common.LowA
	}
	if strings.HasPrefix(*note, "B") {
		return common.B
	}
	if strings.HasPrefix(*note, "C") {
		return common.C
	}
	if strings.HasPrefix(*note, "D") {
		return common.D
	}
	if strings.HasPrefix(*note, "E") {
		return common.E
	}
	if strings.HasPrefix(*note, "F") {
		return common.F
	}
	if strings.HasPrefix(*note, "HG") {
		return common.HighG
	}
	if strings.HasPrefix(*note, "HA") {
		return common.HighA
	}

	return common.NoPitch
}
func lengthFromSuffix(note *string) common.Length {
	if strings.HasSuffix(*note, "16") {
		return common.Sixteenth
	}
	if strings.HasSuffix(*note, "32") {
		return common.Thirtysecond
	}
	if strings.HasSuffix(*note, "1") {
		return common.Whole
	}
	if strings.HasSuffix(*note, "2") {
		return common.Half
	}
	if strings.HasSuffix(*note, "4") {
		return common.Quarter
	}
	if strings.HasSuffix(*note, "8") {
		return common.Eighth
	}

	return common.NoLength
}
