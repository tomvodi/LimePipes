package bww

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/barline"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/boundary"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/length"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/measure"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/pitch"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/symbols/accidental"
	emb "github.com/tomvodi/limepipes-music-model/musicmodel/v1/symbols/embellishment"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/symbols/movement"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/symbols/tie"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/symbols/timeline"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/symbols/tuplet"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/tune"
	"github.com/tomvodi/limepipes/internal/common/music_model"
	"reflect"
	"strconv"
	"strings"
)

type staffContext struct {
	PendingOldTie         pitch.Pitch
	PendingNewTie         bool
	NextMeasureComments   []string
	NextMeasureInLineText []string
	PreviousStaveMeasures []*measure.Measure
}

func convertGrammarToModel(grammar *BwwDocument) (music_model.MusicModel, error) {
	var tunes music_model.MusicModel

	var newTune *tune.Tune
	staffCtx := &staffContext{
		PendingOldTie: pitch.Pitch_NoPitch,
	}
	for i, t := range grammar.Tunes {
		if t.Header.HasTitle() {
			if newTune != nil {
				tunes = append(tunes, newTune)
			}
			newTune = &tune.Tune{}
			if err := fillTuneWithParameter(newTune, t.Header.TuneParameter); err != nil {
				return nil, err
			}
		} else {
			if i == 0 && newTune == nil {
				log.Warn().Msgf("first t doesn't have a title. Setting it to 'no title'")
				newTune = &tune.Tune{}
				if err := fillTuneWithParameter(newTune, t.Header.TuneParameter); err != nil {
					return nil, err
				}
				newTune.Title = "no title"
			} else {
				staffCtx.NextMeasureComments = t.Header.GetComments()
				staffCtx.NextMeasureInLineText = t.Header.GetInlineTexts()
			}
		}

		// TODO when tempo only of first t, set to other tunes as well?
		if err := fillTunePartsFromStaves(newTune, t.Body.Staffs, staffCtx); err != nil {
			return nil, err
		}
	}
	tunes = append(tunes, newTune)

	return tunes, nil
}

func fillTuneWithParameter(tune *tune.Tune, params []*TuneParameter) error {
	for _, param := range params {
		if param.Tempo != nil {
			tempo, err := strconv.ParseUint(param.Tempo.Tempo, 10, 64)
			if err != nil {
				return fmt.Errorf("failed parsing tune tempo: %s", err.Error())
			}

			tune.Tempo = uint32(tempo)
		}

		if param.Description != nil {
			firstParam := param.Description.ParamList[0]
			text := param.Description.Text
			if firstParam == TitleParameter {
				tune.Title = text
			}
			if firstParam == TypeParameter {
				tune.Type = text
			}
			if firstParam == ComposerParameter {
				tune.Composer = text
			}
			if firstParam == FooterParameter {
				tune.Footer = append(tune.Footer, text)
			}
			if firstParam == InlineParameter {
				tune.InLineText = append(tune.InLineText, text)
			}
		}

		if param.Comment != "" {
			tune.Comments = append(tune.Comments, param.Comment)
		}
	}

	return nil
}

func fillTunePartsFromStaves(
	tune *tune.Tune,
	staves []*Staff,
	staffCtx *staffContext,
) error {
	var measures []*measure.Measure

	for _, stave := range staves {
		staveMeasures, err := getMeasuresFromStave(stave, staffCtx)
		if err != nil {
			return err
		}
		staffCtx.PreviousStaveMeasures = staveMeasures

		measures = append(measures, staveMeasures...)
		if stave.End != "" &&
			stave.Dalsegno == nil &&
			stave.DacapoAlFine == nil &&
			len(measures) >= 1 {
			lastMeasure := measures[len(measures)-1]
			var rightBarline *barline.Barline
			if stave.End == "''!I" {
				rightBarline = &barline.Barline{
					Type: barline.Type_Heavy,
					Time: barline.Time_Repeat,
				}
			}
			if stave.End == "!I" {
				rightBarline = &barline.Barline{
					Type: barline.Type_Heavy,
				}
			}
			lastMeasure.RightBarline = rightBarline
		}
	}
	tune.Measures = append(tune.Measures, measures...)

	return nil
}

func getMeasuresFromStave(stave *Staff, ctx *staffContext) ([]*measure.Measure, error) {
	var measures []*measure.Measure
	currMeasure := &measure.Measure{}
	currMeasure.InlineText = ctx.NextMeasureInLineText
	currMeasure.Comments = ctx.NextMeasureComments
	ctx.NextMeasureInLineText = nil
	ctx.NextMeasureComments = nil

	for _, staffSym := range stave.Symbols {
		// if staffSym bar or part start => new measure
		// currMeasure to return measures
		if staffSym.Barline != nil ||
			staffSym.PartStart != nil ||
			staffSym.NextStaffStart != nil {
			measures = cleanupAndAppendMeasure(measures, currMeasure)

			var leftBarline *barline.Barline
			if staffSym.PartStart != nil {
				leftBarline = &barline.Barline{
					Type: barline.Type_Heavy,
				}
				ps := *staffSym.PartStart
				if ps == "I!''" {
					leftBarline.Time = barline.Time_Repeat
				}
			}

			currMeasure = &measure.Measure{
				LeftBarline: leftBarline,
			}

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

		var lastSym *symbols.Symbol
		measSymLen := len(currMeasure.Symbols)
		if len(currMeasure.Symbols) > 0 {
			lastSym = currMeasure.Symbols[measSymLen-1]
		}

		if staffSym.TieOld != nil {
			err := appendTieStartToPreviousNote(*staffSym.TieOld, lastSym, measures, currMeasure, ctx)
			if err != nil {
				return nil, err
			}
		}

		newSym, err := appendStaffSymbolToMeasureSymbols(staffSym, lastSym, currMeasure, measures, ctx)
		if err != nil {
			return nil, err
		}
		if newSym != nil {
			if ctx.PendingOldTie != pitch.Pitch_NoPitch {
				if newSym.Note == nil || !newSym.Note.IsValid() {
					log.Error().Msgf("old tie on pitch %s was started in previous measure but there is "+
						"no note at the beginning of new measure", ctx.PendingOldTie.String())
				} else {
					newSym.Note.Tie = tie.Tie_End
				}

				ctx.PendingOldTie = pitch.Pitch_NoPitch
			}
			currMeasure.Symbols = append(currMeasure.Symbols, newSym)
		}
	}

	// append time line end to measure symbols if it was outside of stave
	if stave.TimelineEnd != nil {
		tl := newTimeLineEnd(stave.TimelineEnd)
		currMeasure.Symbols = append(currMeasure.Symbols, tl)
	}
	if stave.Dalsegno != nil {
		currMeasure.RightBarline = &barline.Barline{
			Type: barline.Type_Regular,
			Time: barline.Time_Dalsegno,
		}
	}
	if stave.DacapoAlFine != nil {
		currMeasure.RightBarline = &barline.Barline{
			Type: barline.Type_Regular,
			Time: barline.Time_DacapoAlFine,
		}
	}

	measures = cleanupAndAppendMeasure(measures, currMeasure)
	return measures, nil
}

func getLastSymbolFromMeasures(measures []*measure.Measure) *symbols.Symbol {
	if len(measures) > 0 {
		lastMeasure := measures[len(measures)-1]
		if len(lastMeasure.Symbols) > 0 {
			return lastMeasure.Symbols[len(lastMeasure.Symbols)-1]
		}
	}

	return nil
}

func handleTriplet(measure *measure.Measure, sym string) error {

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
			symBeforeTriplet.Tuplet.BoundaryType == boundary.Boundary_Start {
			tupletHasAlreadyAStartSymbol = true
		}
	}
	if !tupletHasAlreadyAStartSymbol {
		tripletStartSym := newIrregularGroup(boundary.Boundary_Start, tuplet.Type32)
		measure.Symbols = append(
			measure.Symbols[:tripletStartIdx+1],
			measure.Symbols[tripletStartIdx:]...,
		)
		measure.Symbols[tripletStartIdx] = tripletStartSym
	}

	measure.Symbols = append(measure.Symbols, newIrregularGroup(boundary.Boundary_End, tuplet.Type32))

	return nil
}

func appendTieStartToPreviousNote(
	staffSym string,
	lastSym *symbols.Symbol,
	measures []*measure.Measure,
	currentMeasure *measure.Measure,
	ctx *staffContext,
) error {
	if lastSym == nil {
		lastSym = getLastSymbolFromMeasures(measures)
		if lastSym == nil && len(ctx.PreviousStaveMeasures) > 0 {
			lastSym = getLastSymbolFromMeasures(ctx.PreviousStaveMeasures)
		}
		if lastSym == nil {
			msg := fmt.Sprintf("tie in old format (%s) must follow a note and can't be the first symbol in a measure", staffSym)
			currentMeasure.AddMessage(&measure.ImportMessage{
				Symbol:   staffSym,
				Severity: measure.Severity_Warning,
				Text:     msg,
				Fix:      measure.Fix_SkipSymbol,
			})
			return nil
		}

		if lastSym.IsValidNote() {
			lastSym.Note.Tie = tie.Tie_Start
		} else {
			return fmt.Errorf("tie in old format (%s) must follow a note", staffSym)
		}
	}
	if lastSym.Note == nil {
		return fmt.Errorf("tie in old format (%s) must follow a sym", staffSym)
	}
	if !lastSym.Note.IsValid() {
		msg := fmt.Sprintf(
			"tie in old format (%s) must follow a note with pitch and length",
			staffSym,
		)
		currentMeasure.AddMessage(&measure.ImportMessage{
			Symbol:   staffSym,
			Severity: measure.Severity_Error,
			Text:     msg,
			Fix:      measure.Fix_SkipSymbol,
		})
		return nil
	}
	lastSym.Note.Tie = tie.Tie_Start
	tiePitch := pitchFromSuffix(staffSym)
	ctx.PendingOldTie = tiePitch

	return nil
}

func cleanupAndAppendMeasure(
	measures []*measure.Measure,
	measure *measure.Measure,
) []*measure.Measure {
	cleanupMeasure(measure)
	if len(measure.Symbols) == 0 {
		measure.Symbols = nil
	}

	if !measure.IsNil() {
		return append(measures, measure)
	}

	return measures
}

// cleanupMeasure removes invalid symbols from the measure
// this may be the case for the accidentals at the beginning of the measure which are
// indicating the key of the measure. For bagpipes the key is always sharpf sharpc,
// so we delete these symbols here.
func cleanupMeasure(meas *measure.Measure) {
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

func setMeasureTimeSig(measure *measure.Measure, timeSigSym *string) {
	timeSig := timeSigFromSymbol(timeSigSym)
	measure.Time = timeSig
}

func timeSigFromSymbol(sym *string) *measure.TimeSignature {
	if *sym == "C" {
		return &measure.TimeSignature{
			Beats:    4,
			BeatType: 4,
		}
	}

	if *sym == "C_" {
		return &measure.TimeSignature{
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

	return &measure.TimeSignature{
		Beats:    uint32(beat),
		BeatType: uint32(beatTime),
	}
}

func symbolIndexOf(symbols []*symbols.Symbol, findSym *symbols.Symbol) int {
	for i, symbol := range symbols {
		if reflect.DeepEqual(symbol, findSym) {
			return i
		}
	}
	return -1
}

func removeSymbol(symbols []*symbols.Symbol, idx int) []*symbols.Symbol {
	return append(symbols[:idx], symbols[idx+1:]...)
}

func appendStaffSymbolToMeasureSymbols(
	staffSym *StaffSymbols,
	lastSym *symbols.Symbol,
	currentMeasure *measure.Measure,
	currentStaffMeasures []*measure.Measure,
	ctx *staffContext,
) (*symbols.Symbol, error) {
	newSym := &symbols.Symbol{}

	if staffSym.WholeNote != nil || staffSym.HalfNote != nil ||
		staffSym.QuarterNote != nil || staffSym.EighthNote != nil ||
		staffSym.SixteenthNote != nil || staffSym.ThirtysecondNote != nil {
		// add melody note to last note if it is an emb
		if lastSym != nil && lastSym.Note != nil && lastSym.Note.IsIncomplete() {
			handleNote(staffSym, lastSym.Note)
			return nil, nil
		} else {
			newSym.Note = &symbols.Note{}
			handleNote(staffSym, newSym.Note)
			return newSym, nil
		}
	}
	if staffSym.Segno != nil {
		if currentMeasure.LeftBarline == nil {
			currentMeasure.LeftBarline = &barline.Barline{
				Type: barline.Type_Regular,
				Time: barline.Time_Segno,
			}
		} else {
			currentMeasure.LeftBarline.Time = barline.Time_Segno
		}
	}
	if staffSym.Fine != nil {
		if currentMeasure.RightBarline == nil {
			currentMeasure.RightBarline = &barline.Barline{
				Type: barline.Type_Regular,
				Time: barline.Time_Fine,
			}
		} else {
			currentMeasure.LeftBarline.Time = barline.Time_Fine
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
			Tie: tie.Tie_Start,
		}
		ctx.PendingNewTie = true
		return newSym, nil
	}
	if staffSym.TieEnd != nil {
		// TODO: check if tie start note has same pitch

		// check if the recognized symbol may be an old style tie on E.
		// if so, add tie start to previous note
		oldTieEndE := "^te"
		if *staffSym.TieEnd == oldTieEndE && !ctx.PendingNewTie {
			err := appendTieStartToPreviousNote(oldTieEndE, lastSym, currentStaffMeasures, currentMeasure, ctx)
			if err != nil {
				return nil, err
			}
		}

		if lastSym != nil && lastSym.Note != nil && ctx.PendingNewTie {
			lastSym.Note.Tie = tie.Tie_End
			ctx.PendingNewTie = false
		}
	}
	if staffSym.Flat != nil {
		return handleAccidential(accidental.Accidental_Flat), nil
	}
	if staffSym.Natural != nil {
		return handleAccidential(accidental.Accidental_Natural), nil
	}
	if staffSym.Sharp != nil {
		return handleAccidential(accidental.Accidental_Sharp), nil
	}
	if staffSym.Doubling != nil {
		return handleEmbellishment(emb.Type_Doubling)
	}
	if staffSym.HalfDoubling != nil {
		return handleVariant(emb.Type_Doubling, emb.Variant_Half, emb.Weight_NoWeight)
	}
	if staffSym.ThumbDoubling != nil {
		return handleVariant(emb.Type_Doubling, emb.Variant_Thumb, emb.Weight_NoWeight)
	}
	if staffSym.Grip != nil {
		return handleEmbellishment(emb.Type_Grip)
	}
	if staffSym.GGrip != nil {
		return handleVariant(emb.Type_Grip, emb.Variant_G, emb.Weight_NoWeight)
	}
	if staffSym.ThumbGrip != nil {
		return handleVariant(emb.Type_Grip, emb.Variant_Thumb, emb.Weight_NoWeight)
	}
	if staffSym.HalfGrip != nil {
		return handleVariant(emb.Type_Grip, emb.Variant_Half, emb.Weight_NoWeight)
	}
	if staffSym.Taorluath != nil {
		return handleEmbellishment(emb.Type_Taorluath)
	}
	if staffSym.Bubbly != nil {
		return handleEmbellishment(emb.Type_Bubbly)
	}
	if staffSym.ThrowD != nil {
		return handleVariant(emb.Type_ThrowD, emb.Variant_NoVariant, emb.Weight_Light)
	}
	if staffSym.HeavyThrowD != nil {
		return handleEmbellishment(emb.Type_ThrowD)
	}
	if staffSym.Birl != nil {
		return handleEmbellishment(emb.Type_Birl)
	}
	if staffSym.ABirl != nil {
		return handleEmbellishment(emb.Type_ABirl)
	}
	if staffSym.Strike != nil {
		return handleEmbellishment(emb.Type_Strike)
	}
	if staffSym.GStrike != nil {
		return handleVariant(emb.Type_Strike, emb.Variant_G, emb.Weight_NoWeight)
	}
	if staffSym.LightGStrike != nil {
		return handleVariant(emb.Type_Strike, emb.Variant_G, emb.Weight_Light)
	}
	if staffSym.LightDoubleStrike != nil {
		return handleVariant(emb.Type_DoubleStrike, emb.Variant_NoVariant, emb.Weight_Light)
	}
	if staffSym.DoubleStrike != nil {
		return handleVariant(emb.Type_DoubleStrike, emb.Variant_NoVariant, emb.Weight_NoWeight)
	}
	if staffSym.LightGDoubleStrike != nil {
		return handleVariant(emb.Type_DoubleStrike, emb.Variant_G, emb.Weight_Light)
	}
	if staffSym.GDoubleStrike != nil {
		return handleVariant(emb.Type_DoubleStrike, emb.Variant_G, emb.Weight_NoWeight)
	}
	if staffSym.LightThumbDoubleStrike != nil {
		return handleVariant(emb.Type_DoubleStrike, emb.Variant_Thumb, emb.Weight_Light)
	}
	if staffSym.ThumbDoubleStrike != nil {
		return handleVariant(emb.Type_DoubleStrike, emb.Variant_Thumb, emb.Weight_NoWeight)
	}
	if staffSym.LightHalfDoubleStrike != nil {
		return handleVariant(emb.Type_DoubleStrike, emb.Variant_Half, emb.Weight_Light)
	}
	if staffSym.HalfDoubleStrike != nil {
		return handleVariant(emb.Type_DoubleStrike, emb.Variant_Half, emb.Weight_NoWeight)
	}
	if staffSym.LightTripleStrike != nil {
		return handleVariant(emb.Type_TripleStrike, emb.Variant_NoVariant, emb.Weight_Light)
	}
	if staffSym.TripleStrike != nil {
		return handleVariant(emb.Type_TripleStrike, emb.Variant_NoVariant, emb.Weight_NoWeight)
	}
	if staffSym.LightGTripleStrike != nil {
		return handleVariant(emb.Type_TripleStrike, emb.Variant_G, emb.Weight_Light)
	}
	if staffSym.GTripleStrike != nil {
		return handleVariant(emb.Type_TripleStrike, emb.Variant_G, emb.Weight_NoWeight)
	}
	if staffSym.LightThumbTripleStrike != nil {
		return handleVariant(emb.Type_TripleStrike, emb.Variant_Thumb, emb.Weight_Light)
	}
	if staffSym.ThumbTripleStrike != nil {
		return handleVariant(emb.Type_TripleStrike, emb.Variant_Thumb, emb.Weight_NoWeight)
	}
	if staffSym.LightHalfTripleStrike != nil {
		return handleVariant(emb.Type_TripleStrike, emb.Variant_Half, emb.Weight_Light)
	}
	if staffSym.HalfTripleStrike != nil {
		return handleVariant(emb.Type_TripleStrike, emb.Variant_Half, emb.Weight_NoWeight)
	}
	if staffSym.DDoubleGrace != nil {
		return handleDoubleGrace(pitch.Pitch_D)
	}
	if staffSym.EDoubleGrace != nil {
		return handleDoubleGrace(pitch.Pitch_E)
	}
	if staffSym.FDoubleGrace != nil {
		return handleDoubleGrace(pitch.Pitch_F)
	}
	if staffSym.GDoubleGrace != nil {
		return handleDoubleGrace(pitch.Pitch_HighG)
	}
	if staffSym.ThumbDoubleGrace != nil {
		return handleDoubleGrace(pitch.Pitch_HighA)
	}
	if staffSym.HalfStrike != nil {
		return handleVariant(emb.Type_Strike, emb.Variant_Half, emb.Weight_NoWeight)
	}
	if staffSym.LightHalfStrike != nil {
		return handleVariant(emb.Type_Strike, emb.Variant_Half, emb.Weight_Light)
	}
	if staffSym.ThumbStrike != nil {
		return handleVariant(emb.Type_Strike, emb.Variant_Thumb, emb.Weight_NoWeight)
	}
	if staffSym.LightThumbStrike != nil {
		return handleVariant(emb.Type_Strike, emb.Variant_Thumb, emb.Weight_Light)
	}
	if staffSym.Pele != nil {
		return handleEmbellishment(emb.Type_Pele)
	}
	if staffSym.LightPele != nil {
		return handleVariant(emb.Type_Pele, emb.Variant_NoVariant, emb.Weight_Light)
	}
	if staffSym.ThumbPele != nil {
		return handleVariant(emb.Type_Pele, emb.Variant_Thumb, emb.Weight_NoWeight)
	}
	if staffSym.LightThumbPele != nil {
		return handleVariant(emb.Type_Pele, emb.Variant_Thumb, emb.Weight_Light)
	}
	if staffSym.HalfPele != nil {
		return handleVariant(emb.Type_Pele, emb.Variant_Half, emb.Weight_NoWeight)
	}
	if staffSym.IrregularGroupStart != nil {
		ttype := tupletTypeFromSymbol(staffSym.IrregularGroupStart)
		return handleIrregularGroup(boundary.Boundary_Start, ttype)
	}
	if staffSym.IrregularGroupEnd != nil {
		ttype := tupletTypeFromSymbol(staffSym.IrregularGroupEnd)
		// handling old style ties on E (^3e) but ignoring an error
		if ttype == tuplet.Type32 {
			_ = handleTriplet(currentMeasure, "^3e")
			return nil, nil
		} else {
			return handleIrregularGroup(boundary.Boundary_End, ttype)
		}
	}
	if staffSym.LightHalfPele != nil {
		return handleVariant(emb.Type_Pele, emb.Variant_Half, emb.Weight_Light)
	}
	if staffSym.GBirl != nil ||
		staffSym.ThumbBirl != nil {
		return handleEmbellishment(emb.Type_GraceBirl)
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
	if staffSym.TimelineStart != nil {
		return handleTimeLine(*staffSym.TimelineStart)
	}
	if staffSym.TimelineEnd != nil {
		return newTimeLineEnd(staffSym.TimelineEnd), nil
	}
	if staffSym.Comment != nil {
		handleInsideStaffComment(lastSym, currentMeasure, *staffSym.Comment)
	}
	if staffSym.Description != nil {
		handleInsideStaffComment(lastSym, currentMeasure, staffSym.Description.Text)
	}
	if staffSym.Cadence != nil {
		return handleCadence(staffSym.Cadence, false)
	}
	if staffSym.FermatCadence != nil {
		return handleCadence(staffSym.FermatCadence, true)
	}
	if staffSym.Embari != nil {
		return handleMovement(movement.Type_Embari, staffSym.Embari, true, true)
	}
	if staffSym.Endari != nil {
		return handleMovement(movement.Type_Endari, staffSym.Endari, true, true)
	}
	if staffSym.Chedari != nil {
		return handleMovement(movement.Type_Chedari, staffSym.Chedari, true, true)
	}
	if staffSym.Hedari != nil {
		return handleMovement(movement.Type_Hedari, staffSym.Hedari, true, false)
	}
	if staffSym.Dili != nil {
		return handleMovement(movement.Type_Dili, staffSym.Dili, true, true)
	}
	if staffSym.Tra != nil {
		return handleMovement(movement.Type_Tra, staffSym.Tra, false, true)
	}
	if staffSym.Edre != nil {
		mv, _ := handleMovement(movement.Type_Edre, staffSym.Edre, true, true)
		p := pitchFromSuffix(*staffSym.Edre)
		if p != pitch.Pitch_E {
			mv.Note.Movement.PitchHint = p
		}
		return mv, nil
	}
	if staffSym.HalfEdre != nil {
		mv, _ := handleMovement(movement.Type_Edre, staffSym.HalfEdre, true, true)
		mv.Note.Movement.Variant = movement.Variant_Half
		return mv, nil
	}
	if staffSym.GEdre != nil {
		return handleMovement(movement.Type_Edre, staffSym.GEdre, true, true)
	}
	if staffSym.ThumbEdre != nil {
		return handleMovement(movement.Type_Edre, staffSym.ThumbEdre, true, true)
	}
	if staffSym.Dare != nil {
		return handleMovement(movement.Type_Dare, staffSym.Dare, true, true)
	}
	if staffSym.HalfDare != nil {
		return handleMovement(movement.Type_Dare, staffSym.HalfDare, true, true)
	}
	if staffSym.ThumbDare != nil {
		return handleMovement(movement.Type_Dare, staffSym.ThumbDare, true, true)
	}
	if staffSym.GDare != nil {
		return handleMovement(movement.Type_Dare, staffSym.GDare, true, true)
	}
	if staffSym.CheCheRe != nil {
		return handleMovement(movement.Type_CheCheRe, staffSym.CheCheRe, true, true)
	}
	if staffSym.HalfCheCheRe != nil {
		return handleMovement(movement.Type_CheCheRe, staffSym.HalfCheCheRe, true, true)
	}
	if staffSym.ThumbCheCheRe != nil {
		return handleMovement(movement.Type_CheCheRe, staffSym.ThumbCheCheRe, true, true)
	}
	if staffSym.GripAbbrev != nil {
		return handleMovement(movement.Type_Grip, staffSym.GripAbbrev, true, true)
	}
	if staffSym.Deda != nil {
		return handleMovement(movement.Type_Deda, staffSym.Deda, true, true)
	}
	if staffSym.Enbain != nil {
		return handleMovement(movement.Type_Enbain, staffSym.Enbain, true, true)
	}
	if staffSym.GEnbain != nil {
		return handleMovement(movement.Type_Enbain, staffSym.GEnbain, true, true)
	}
	if staffSym.ThumbEnbain != nil {
		return handleMovement(movement.Type_Enbain, staffSym.ThumbEnbain, true, true)
	}
	if staffSym.Otro != nil {
		return handleMovement(movement.Type_Otro, staffSym.Otro, true, true)
	}
	if staffSym.GOtro != nil {
		return handleMovement(movement.Type_Otro, staffSym.GOtro, true, true)
	}
	if staffSym.ThumbOtro != nil {
		return handleMovement(movement.Type_Otro, staffSym.ThumbOtro, true, true)
	}
	if staffSym.Odro != nil {
		return handleMovement(movement.Type_Odro, staffSym.Odro, true, true)
	}
	if staffSym.GOdro != nil {
		return handleMovement(movement.Type_Odro, staffSym.GOdro, true, true)
	}
	if staffSym.ThumbOdro != nil {
		return handleMovement(movement.Type_Odro, staffSym.ThumbOdro, true, true)
	}
	if staffSym.Adeda != nil {
		return handleMovement(movement.Type_Adeda, staffSym.Adeda, true, true)
	}
	if staffSym.GAdeda != nil {
		return handleMovement(movement.Type_Adeda, staffSym.GAdeda, true, true)
	}
	if staffSym.ThumbAdeda != nil {
		return handleMovement(movement.Type_Adeda, staffSym.ThumbAdeda, true, true)
	}
	if staffSym.EchoBeats != nil {
		mv, _ := handleMovement(movement.Type_EchoBeat, staffSym.EchoBeats, false, false)
		pitch := pitchFromSuffix(*staffSym.EchoBeats)
		mv.Note.Movement.Pitch = pitch
		return mv, nil
	}
	if staffSym.Darodo != nil {
		return handleMovement(movement.Type_Darodo, staffSym.Darodo, false, true)
	}
	if staffSym.Hiharin != nil {
		return handleMovement(movement.Type_Hiharin, staffSym.Hiharin, false, false)
	}
	if staffSym.Rodin != nil {
		return handleMovement(movement.Type_Rodin, staffSym.Rodin, false, false)
	}
	if staffSym.Chelalho != nil {
		return handleMovement(movement.Type_Chelalho, staffSym.Chelalho, false, false)
	}
	if staffSym.Din != nil {
		return handleMovement(movement.Type_Din, staffSym.Din, false, false)
	}
	if staffSym.Lemluath != nil {
		lemSym, hadBrea := stripBreabach(staffSym.Lemluath)
		mv, _ := handleMovement(movement.Type_Lemluath, &lemSym, false, true)
		pitch := pitchFromSuffix(lemSym)
		mv.Note.Movement.PitchHint = pitch
		mv.Note.Movement.Breabach = hadBrea
		return mv, nil
	}
	if staffSym.LemluathAbbrev != nil {
		lemSym, hadBrea := stripBreabach(staffSym.LemluathAbbrev)
		if lastSym == nil || lastSym.Note == nil {
			return nil, fmt.Errorf("lemluath abbreviation %s must follow melody note", *staffSym.LemluathAbbrev)
		}
		if !lastSym.Note.IsValid() {
			return nil, fmt.Errorf("lemluath abbreviation %s must follow a valid melody note", *staffSym.LemluathAbbrev)
		}
		sym, _ := handleMovement(movement.Type_Lemluath, &lemSym, false, true)
		move := sym.Note.Movement
		pitch := pitchFromSuffix(lemSym)
		move.PitchHint = pitch
		move.Breabach = hadBrea
		lastSym.Note.Movement = move
		return nil, nil
	}
	if staffSym.TaorluathPio != nil {
		return handleMovementWithPitchHintSuffixAndBreabach(
			movement.Type_Taorluath, staffSym.TaorluathPio, false, true,
		)
	}
	if staffSym.TaorluathAbbrev != nil {
		if lastSym == nil || lastSym.Note == nil {
			return nil, fmt.Errorf("taorluath abbreviation %s must follow melody note", *staffSym.TaorluathAbbrev)
		}
		if !lastSym.Note.IsValid() {
			return nil, fmt.Errorf("taorluath abbreviation %s must follow a valid melody note", *staffSym.TaorluathAbbrev)
		}
		sym, err := handleMovementWithPitchHintSuffixAndBreabach(
			movement.Type_Taorluath, staffSym.TaorluathAbbrev, false, true,
		)
		move := sym.Note.Movement
		lastSym.Note.Movement = move
		return nil, err
	}
	if staffSym.TaorluathAmach != nil {
		pitch := pitchFromSuffix(*staffSym.TaorluathAmach)
		mv, _ := handleMovement(movement.Type_Taorluath, staffSym.TaorluathAmach, false, false)
		mv.Note.Movement.AMach = true
		mv.Note.Movement.PitchHint = pitch
		return mv, nil
	}
	if staffSym.Crunluath != nil {
		mv, _ := handleMovementWithPitchHintSuffixAndBreabach(
			movement.Type_Crunluath, staffSym.Crunluath, false, true,
		)
		if strings.Contains(*staffSym.Crunluath, "crunllgla") {
			mv.Note.Movement.AdditionalPitchHint = pitch.Pitch_LowG
		}
		return mv, nil
	}
	if staffSym.CrunluathAbbrev != nil {
		if lastSym == nil || lastSym.Note == nil {
			return nil, fmt.Errorf("crunluath abbreviation %s must follow melody note", *staffSym.CrunluathAbbrev)
		}
		if !lastSym.Note.IsValid() {
			return nil, fmt.Errorf("crunluath abbreviation %s must follow a valid melody note", *staffSym.CrunluathAbbrev)
		}
		sym, err := handleMovementWithPitchHintSuffixAndBreabach(
			movement.Type_Crunluath, staffSym.CrunluathAbbrev, false, true,
		)
		move := sym.Note.Movement
		lastSym.Note.Movement = move
		return nil, err
	}
	if staffSym.CrunluathAmach != nil {
		pitch := pitchFromSuffix(*staffSym.CrunluathAmach)
		mv, _ := handleMovement(movement.Type_Crunluath, staffSym.CrunluathAmach, false, false)
		mv.Note.Movement.AMach = true
		mv.Note.Movement.PitchHint = pitch
		return mv, nil
	}
	if staffSym.Tripling != nil {
		pitch := pitchFromSuffix(*staffSym.Tripling)
		mv, _ := handleMovement(movement.Type_Tripling, staffSym.Tripling, false, true)
		// thumb variant handled here because pt is always recognized as thumb
		if strings.HasPrefix(*staffSym.Tripling, "ptt") {
			mv.Note.Movement.Variant = movement.Variant_Thumb
		}
		mv.Note.Movement.Pitch = pitch
		return mv, nil
	}
	if staffSym.Tempo != nil {
		tempo, err := strconv.ParseUint(staffSym.Tempo.Tempo, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed parsing tune tempo: %s", err.Error())
		}
		return &symbols.Symbol{TempoChange: &tempo}, nil
	}

	return nil, nil // fmt.Errorf("staff symbol %v not handled", staffSym)
}

func handleMovementWithPitchHintSuffixAndBreabach(
	mtype movement.Type,
	sym *string,
	withThumb bool,
	withHalf bool,
) (*symbols.Symbol, error) {
	strippedSym, hadBrea := stripBreabach(sym)
	currSym, _ := handleMovement(mtype, &strippedSym, withThumb, withHalf)
	pitch := pitchFromSuffix(strippedSym)
	currSym.Note.Movement.PitchHint = pitch
	currSym.Note.Movement.Breabach = hadBrea

	return currSym, nil
}

func handleMovement(mtype movement.Type, sym *string, withThumb bool, withHalf bool) (*symbols.Symbol, error) {
	showAbbr := false
	if strings.HasPrefix(*sym, "p") {
		showAbbr = true
	}
	mVar := movement.Variant_NoVariant
	if withHalf {
		if strings.HasPrefix(*sym, "h") || strings.HasPrefix(*sym, "ph") {
			mVar = movement.Variant_Half
		}
	}

	if withThumb {
		if strings.HasPrefix(*sym, "t") || strings.HasPrefix(*sym, "pt") {
			mVar = movement.Variant_Thumb
		}
	}

	if strings.HasPrefix(*sym, "g") {
		mVar = movement.Variant_G
	}

	if strings.HasSuffix(*sym, "8") ||
		strings.HasSuffix(*sym, "16") {
		mVar = movement.Variant_LongLowG
	}

	return &symbols.Symbol{
		Note: &symbols.Note{
			Movement: &movement.Movement{
				Type:       mtype,
				Abbreviate: showAbbr,
				Variant:    mVar,
			},
		},
	}, nil
}

func handleCadence(
	cad *string,
	fermata bool,
) (*symbols.Symbol, error) {
	return &symbols.Symbol{
		Note: &symbols.Note{
			Movement: &movement.Movement{
				Type:    movement.Type_Cadence,
				Fermata: fermata,
				Pitches: pitchesFromCadenceSym(*cad, fermata),
			},
		},
	}, nil
}

func pitchesFromCadenceSym(sym string, fermata bool) []pitch.Pitch {
	if fermata {
		sym = strings.Replace(sym, "fcad", "", 1)
	} else {
		sym = strings.Replace(sym, "cad", "", 1)
	}

	pitches := make([]pitch.Pitch, len(sym))
	for i, ch := range sym {
		switch ch {
		case 'g':
			pitches[i] = pitch.Pitch_HighG
		case 'e':
			pitches[i] = pitch.Pitch_E
		case 'd':
			pitches[i] = pitch.Pitch_D
		case 'a':
			pitches[i] = pitch.Pitch_HighA
		case 'f':
			pitches[i] = pitch.Pitch_F
		default:
			log.Error().Msgf("char %c is not handled for cadence symbol", ch)
			pitches[i] = pitch.Pitch_NoPitch
		}
	}
	return pitches
}

func handleInsideStaffComment(
	lastSym *symbols.Symbol,
	currentMeasure *measure.Measure,
	text string,
) {
	if lastSym != nil {
		if lastSym.IsNote() {
			lastSym.Note.Comment = text
		}
	} else {
		if currentMeasure != nil {
			currentMeasure.Comments = append(currentMeasure.Comments, text)
		}
	}
}

func handleEmbellishment(
	embType emb.Type,
) (*symbols.Symbol, error) {
	return &symbols.Symbol{
		Note: &symbols.Note{
			Embellishment: &emb.Embellishment{
				Type: embType,
			},
		},
	}, nil
}

func handleDoubleGrace(pitch pitch.Pitch) (*symbols.Symbol, error) {
	doubleG, err := handleEmbellishment(emb.Type_DoubleGrace)
	doubleG.Note.Embellishment.Pitch = pitch
	return doubleG, err
}

func handleVariant(
	embType emb.Type,
	variant emb.Variant,
	weight emb.Weight,
) (*symbols.Symbol, error) {
	return &symbols.Symbol{
		Note: &symbols.Note{
			Embellishment: &emb.Embellishment{
				Type:    embType,
				Variant: variant,
				Weight:  weight,
			},
		},
	}, nil
}

func handleIrregularGroup(
	boundary boundary.Boundary,
	ttype tuplet.TupletType,
) (*symbols.Symbol, error) {
	return newIrregularGroup(boundary, ttype), nil
}

func newIrregularGroup(boundary boundary.Boundary,
	ttype tuplet.TupletType,
) *symbols.Symbol {
	tpl := tuplet.NewTuplet(boundary, ttype)
	return &symbols.Symbol{
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

func handleDots(staffSym *StaffSymbols, lastSym *symbols.Symbol) {
	var dotCount = uint32(0)
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

func handleAccidential(acc accidental.Accidental) *symbols.Symbol {
	return &symbols.Symbol{
		Note: &symbols.Note{
			Accidental: acc,
		},
	}
}

func handleTimeLine(sym string) (*symbols.Symbol, error) {
	if sym == "'1" {
		return newTimeLineStartSymbol(timeline.Type_First), nil
	}
	if sym == "'2" {
		return newTimeLineStartSymbol(timeline.Type_Second), nil
	}
	if sym == "'22" {
		return newTimeLineStartSymbol(timeline.Type_SecondOf2), nil
	}
	if sym == "'23" {
		return newTimeLineStartSymbol(timeline.Type_SecondOf3), nil
	}
	if sym == "'24" {
		return newTimeLineStartSymbol(timeline.Type_SecondOf4), nil
	}
	if sym == "'224" {
		return newTimeLineStartSymbol(timeline.Type_SecondOf2And4), nil
	}
	if sym == "'25" {
		return newTimeLineStartSymbol(timeline.Type_SecondOf5), nil
	}
	if sym == "'26" {
		return newTimeLineStartSymbol(timeline.Type_SecondOf6), nil
	}
	if sym == "'27" {
		return newTimeLineStartSymbol(timeline.Type_SecondOf7), nil
	}
	if sym == "'28" {
		return newTimeLineStartSymbol(timeline.Type_SecondOf8), nil
	}
	if sym == "'si" {
		return newTimeLineStartSymbol(timeline.Type_Singling), nil
	}
	if sym == "'do" {
		return newTimeLineStartSymbol(timeline.Type_Doubling), nil
	}
	if sym == "'bis" {
		return newTimeLineStartSymbol(timeline.Type_Bis), nil
	}
	if sym == "'intro" {
		return newTimeLineStartSymbol(timeline.Type_Intro), nil
	}

	return nil, fmt.Errorf("time line symbol %s not handled", sym)
}

func newTimeLineStartSymbol(ttype timeline.Type) *symbols.Symbol {
	return &symbols.Symbol{
		Timeline: &timeline.TimeLine{
			BoundaryType: boundary.Boundary_Start,
			Type:         ttype,
		},
	}
}

func embellishmentForSingleGrace(grace *string) *emb.Embellishment {
	emb := &emb.Embellishment{
		Type: emb.Type_SingleGrace,
	}

	if *grace == "ag" {
		emb.Pitch = pitch.Pitch_LowA
	}
	if *grace == "bg" {
		emb.Pitch = pitch.Pitch_B
	}
	if *grace == "cg" {
		emb.Pitch = pitch.Pitch_C
	}
	if *grace == "dg" {
		emb.Pitch = pitch.Pitch_D
	}
	if *grace == "eg" {
		emb.Pitch = pitch.Pitch_E
	}
	if *grace == "fg" {
		emb.Pitch = pitch.Pitch_F
	}
	if *grace == "gg" {
		emb.Pitch = pitch.Pitch_HighG
	}
	if *grace == "tg" {
		emb.Pitch = pitch.Pitch_HighA
	}
	return emb
}

func pitchFromStaffNotePrefix(note *string) pitch.Pitch {
	if strings.HasPrefix(*note, "LG") {
		return pitch.Pitch_LowG
	}
	if strings.HasPrefix(*note, "LA") {
		return pitch.Pitch_LowA
	}
	if strings.HasPrefix(*note, "B") {
		return pitch.Pitch_B
	}
	if strings.HasPrefix(*note, "C") {
		return pitch.Pitch_C
	}
	if strings.HasPrefix(*note, "D") {
		return pitch.Pitch_D
	}
	if strings.HasPrefix(*note, "E") {
		return pitch.Pitch_E
	}
	if strings.HasPrefix(*note, "F") {
		return pitch.Pitch_F
	}
	if strings.HasPrefix(*note, "HG") {
		return pitch.Pitch_HighG
	}
	if strings.HasPrefix(*note, "HA") {
		return pitch.Pitch_HighA
	}

	return pitch.Pitch_NoPitch
}
func lengthFromSuffix(note *string) length.Length {
	if strings.HasSuffix(*note, "16") {
		return length.Length_Sixteenth
	}
	if strings.HasSuffix(*note, "32") {
		return length.Length_Thirtysecond
	}
	if strings.HasSuffix(*note, "1") {
		return length.Length_Whole
	}
	if strings.HasSuffix(*note, "2") {
		return length.Length_Half
	}
	if strings.HasSuffix(*note, "4") {
		return length.Length_Quarter
	}
	if strings.HasSuffix(*note, "8") {
		return length.Length_Eighth
	}

	return length.Length_NoLength
}

func pitchFromSuffix(sym string) pitch.Pitch {
	if strings.HasSuffix(sym, "lg") {
		return pitch.Pitch_LowG
	}
	if strings.HasSuffix(sym, "la") {
		return pitch.Pitch_LowA
	}
	if strings.HasSuffix(sym, "b") {
		return pitch.Pitch_B
	}
	if strings.HasSuffix(sym, "c") {
		return pitch.Pitch_C
	}
	if strings.HasSuffix(sym, "d") {
		return pitch.Pitch_D
	}
	if strings.HasSuffix(sym, "e") {
		return pitch.Pitch_E
	}
	if strings.HasSuffix(sym, "f") {
		return pitch.Pitch_F
	}
	if strings.HasSuffix(sym, "hg") {
		return pitch.Pitch_HighG
	}
	if strings.HasSuffix(sym, "ha") {
		return pitch.Pitch_HighA
	}
	return pitch.Pitch_NoPitch
}

func stripBreabach(sym *string) (string, bool) {
	stripped := strings.Replace(*sym, "brea", "", 1)
	didReplace := len(*sym) != len(stripped)
	return stripped, didReplace
}

func newTimeLineEnd(sym *string) *symbols.Symbol {
	ttype := timeline.Type_NoType
	if *sym == "bis_'" {
		ttype = timeline.Type_Bis
	}
	return &symbols.Symbol{
		Timeline: &timeline.TimeLine{
			BoundaryType: boundary.Boundary_End,
			Type:         ttype,
		},
	}
}
