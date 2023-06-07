package bww

import (
	"banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/barline"
	"banduslib/internal/common/music_model/import_message"
	"banduslib/internal/common/music_model/symbols"
	"banduslib/internal/common/music_model/symbols/accidental"
	emb "banduslib/internal/common/music_model/symbols/embellishment"
	"banduslib/internal/common/music_model/symbols/movement"
	"banduslib/internal/common/music_model/symbols/tie"
	"banduslib/internal/common/music_model/symbols/time_line"
	"banduslib/internal/common/music_model/symbols/tuplet"
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"strconv"
	"strings"
)

type staffContext struct {
	PendingOldTie         common.Pitch
	PendingNewTie         bool
	NextMeasureComments   []string
	NextMeasureInLineText []string
	PreviousStaveMeasures []*music_model.Measure
}

func convertGrammarToModel(grammar *BwwDocument) (music_model.MusicModel, error) {
	var tunes music_model.MusicModel

	var newTune *music_model.Tune
	staffCtx := &staffContext{
		PendingOldTie: common.NoPitch,
	}
	for i, tune := range grammar.Tunes {
		if tune.Header.HasTitle() {
			if newTune != nil {
				tunes = append(tunes, newTune)
			}
			newTune = &music_model.Tune{}
			if err := fillTuneWithParameter(newTune, tune.Header.TuneParameter); err != nil {
				return nil, err
			}
		} else {
			if i == 0 && newTune == nil {
				log.Warn().Msgf("first tune doesn't have a title. Setting it to 'no title'")
				newTune = &music_model.Tune{}
				if err := fillTuneWithParameter(newTune, tune.Header.TuneParameter); err != nil {
					return nil, err
				}
				newTune.Title = "no title"
			} else {
				staffCtx.NextMeasureComments = tune.Header.GetComments()
				staffCtx.NextMeasureInLineText = tune.Header.GetInlineTexts()
			}
		}

		// TODO when tempo only of first tune, set to other tunes as well?
		if err := fillTunePartsFromStaves(newTune, tune.Body.Staffs, staffCtx); err != nil {
			return nil, err
		}
	}
	tunes = append(tunes, newTune)

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
	tune *music_model.Tune,
	staves []*Staff,
	staffCtx *staffContext,
) error {
	var measures []*music_model.Measure

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
					Type:     barline.Heavy,
					Timeline: barline.Repeat,
				}
			}
			if stave.End == "!I" {
				rightBarline = &barline.Barline{
					Type: barline.Heavy,
				}
			}
			lastMeasure.RightBarline = rightBarline
		}
	}
	tune.Measures = append(tune.Measures, measures...)

	return nil
}

func getMeasuresFromStave(stave *Staff, ctx *staffContext) ([]*music_model.Measure, error) {
	var measures []*music_model.Measure
	currMeasure := &music_model.Measure{}
	currMeasure.InLineText = ctx.NextMeasureInLineText
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
					Type: barline.Heavy,
				}
				ps := *staffSym.PartStart
				if ps == "I!''" {
					leftBarline.Timeline = barline.Repeat
				}
			}

			currMeasure = &music_model.Measure{
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

		var lastSym *music_model.Symbol
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
			if ctx.PendingOldTie != common.NoPitch {
				if newSym.Note == nil || !newSym.Note.IsValid() {
					log.Error().Msgf("old tie on pitch %s was started in previous measure but there is "+
						"no note at the beginning of new measure", ctx.PendingOldTie.String())
				} else {
					newSym.Note.Tie = tie.End
				}

				ctx.PendingOldTie = common.NoPitch
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
			Type:     barline.Regular,
			Timeline: barline.Dalsegno,
		}
	}
	if stave.DacapoAlFine != nil {
		currMeasure.RightBarline = &barline.Barline{
			Type:     barline.Regular,
			Timeline: barline.DacapoAlFine,
		}
	}

	measures = cleanupAndAppendMeasure(measures, currMeasure)
	return measures, nil
}

func getLastSymbolFromMeasures(measures []*music_model.Measure) *music_model.Symbol {
	if len(measures) > 0 {
		lastMeasure := measures[len(measures)-1]
		if len(lastMeasure.Symbols) > 0 {
			return lastMeasure.Symbols[len(lastMeasure.Symbols)-1]
		}
	}

	return nil
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

func appendTieStartToPreviousNote(
	staffSym string,
	lastSym *music_model.Symbol,
	measures []*music_model.Measure,
	currentMeasure *music_model.Measure,
	ctx *staffContext,
) error {
	if lastSym == nil {
		lastSym = getLastSymbolFromMeasures(measures)
		if lastSym == nil && len(ctx.PreviousStaveMeasures) > 0 {
			lastSym = getLastSymbolFromMeasures(ctx.PreviousStaveMeasures)
		}
		if lastSym == nil {
			msg := fmt.Sprintf("tie in old format (%s) must follow a note and can't be the first symbol in a measure", staffSym)
			currentMeasure.AddMessage(&import_message.ImportMessage{
				Symbol: staffSym,
				Type:   import_message.Warning,
				Text:   msg,
				Fix:    import_message.SkipSymbol,
			})
			return nil
		}

		if lastSym.IsValidNote() {
			lastSym.Note.Tie = tie.Start
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
		currentMeasure.AddMessage(&import_message.ImportMessage{
			Symbol: staffSym,
			Type:   import_message.Error,
			Text:   msg,
			Fix:    import_message.SkipSymbol,
		})
		return nil
	}
	lastSym.Note.Tie = tie.Start
	tiePitch := pitchFromSuffix(staffSym)
	ctx.PendingOldTie = tiePitch

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

	if !measure.IsNil() {
		return append(measures, measure)
	}

	return measures
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
	currentStaffMeasures []*music_model.Measure,
	ctx *staffContext,
) (*music_model.Symbol, error) {
	newSym := &music_model.Symbol{}

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
				Type:     barline.Regular,
				Timeline: barline.Segno,
			}
		} else {
			currentMeasure.LeftBarline.Timeline = barline.Segno
		}
	}
	if staffSym.Fine != nil {
		if currentMeasure.RightBarline == nil {
			currentMeasure.RightBarline = &barline.Barline{
				Type:     barline.Regular,
				Timeline: barline.Fine,
			}
		} else {
			currentMeasure.LeftBarline.Timeline = barline.Fine
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
			Tie: tie.Start,
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
			lastSym.Note.Tie = tie.End
			ctx.PendingNewTie = false
		}
	}
	if staffSym.Flat != nil {
		return handleAccidential(accidental.Flat), nil
	}
	if staffSym.Natural != nil {
		return handleAccidential(accidental.Natural), nil
	}
	if staffSym.Sharp != nil {
		return handleAccidential(accidental.Sharp), nil
	}
	if staffSym.Doubling != nil {
		return handleEmbellishment(emb.Doubling)
	}
	if staffSym.HalfDoubling != nil {
		return handleEmbellishmentVariant(emb.Doubling, emb.Half, emb.NoWeight)
	}
	if staffSym.ThumbDoubling != nil {
		return handleEmbellishmentVariant(emb.Doubling, emb.Thumb, emb.NoWeight)
	}
	if staffSym.Grip != nil {
		return handleEmbellishment(emb.Grip)
	}
	if staffSym.GGrip != nil {
		return handleEmbellishmentVariant(emb.Grip, emb.G, emb.NoWeight)
	}
	if staffSym.ThumbGrip != nil {
		return handleEmbellishmentVariant(emb.Grip, emb.Thumb, emb.NoWeight)
	}
	if staffSym.HalfGrip != nil {
		return handleEmbellishmentVariant(emb.Grip, emb.Half, emb.NoWeight)
	}
	if staffSym.Taorluath != nil {
		return handleEmbellishment(emb.Taorluath)
	}
	if staffSym.Bubbly != nil {
		return handleEmbellishment(emb.Bubbly)
	}
	if staffSym.ThrowD != nil {
		return handleEmbellishmentVariant(emb.ThrowD, emb.NoVariant, emb.Light)
	}
	if staffSym.HeavyThrowD != nil {
		return handleEmbellishment(emb.ThrowD)
	}
	if staffSym.Birl != nil {
		return handleEmbellishment(emb.Birl)
	}
	if staffSym.ABirl != nil {
		return handleEmbellishment(emb.ABirl)
	}
	if staffSym.Strike != nil {
		sym, err := handleEmbellishment(emb.Strike)
		pitch := pitchFromSuffix(*staffSym.Strike)
		sym.Note.Embellishment.Pitch = pitch
		return sym, err
	}
	if staffSym.GStrike != nil {
		return handleEmbellishmentVariant(emb.Strike, emb.G, emb.NoWeight)
	}
	if staffSym.LightGStrike != nil {
		return handleEmbellishmentVariant(emb.Strike, emb.G, emb.Light)
	}
	if staffSym.LightDoubleStrike != nil {
		return handleEmbellishmentVariant(emb.DoubleStrike, emb.NoVariant, emb.Light)
	}
	if staffSym.DoubleStrike != nil {
		return handleEmbellishmentVariant(emb.DoubleStrike, emb.NoVariant, emb.NoWeight)
	}
	if staffSym.LightGDoubleStrike != nil {
		return handleEmbellishmentVariant(emb.DoubleStrike, emb.G, emb.Light)
	}
	if staffSym.GDoubleStrike != nil {
		return handleEmbellishmentVariant(emb.DoubleStrike, emb.G, emb.NoWeight)
	}
	if staffSym.LightThumbDoubleStrike != nil {
		return handleEmbellishmentVariant(emb.DoubleStrike, emb.Thumb, emb.Light)
	}
	if staffSym.ThumbDoubleStrike != nil {
		return handleEmbellishmentVariant(emb.DoubleStrike, emb.Thumb, emb.NoWeight)
	}
	if staffSym.LightHalfDoubleStrike != nil {
		return handleEmbellishmentVariant(emb.DoubleStrike, emb.Half, emb.Light)
	}
	if staffSym.HalfDoubleStrike != nil {
		return handleEmbellishmentVariant(emb.DoubleStrike, emb.Half, emb.NoWeight)
	}
	if staffSym.LightTripleStrike != nil {
		return handleEmbellishmentVariant(emb.TripleStrike, emb.NoVariant, emb.Light)
	}
	if staffSym.TripleStrike != nil {
		return handleEmbellishmentVariant(emb.TripleStrike, emb.NoVariant, emb.NoWeight)
	}
	if staffSym.LightGTripleStrike != nil {
		return handleEmbellishmentVariant(emb.TripleStrike, emb.G, emb.Light)
	}
	if staffSym.GTripleStrike != nil {
		return handleEmbellishmentVariant(emb.TripleStrike, emb.G, emb.NoWeight)
	}
	if staffSym.LightThumbTripleStrike != nil {
		return handleEmbellishmentVariant(emb.TripleStrike, emb.Thumb, emb.Light)
	}
	if staffSym.ThumbTripleStrike != nil {
		return handleEmbellishmentVariant(emb.TripleStrike, emb.Thumb, emb.NoWeight)
	}
	if staffSym.LightHalfTripleStrike != nil {
		return handleEmbellishmentVariant(emb.TripleStrike, emb.Half, emb.Light)
	}
	if staffSym.HalfTripleStrike != nil {
		return handleEmbellishmentVariant(emb.TripleStrike, emb.Half, emb.NoWeight)
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
		return handleEmbellishmentVariant(emb.Strike, emb.Half, emb.NoWeight)
	}
	if staffSym.LightHalfStrike != nil {
		return handleEmbellishmentVariant(emb.Strike, emb.Half, emb.Light)
	}
	if staffSym.ThumbStrike != nil {
		return handleEmbellishmentVariant(emb.Strike, emb.Thumb, emb.NoWeight)
	}
	if staffSym.LightThumbStrike != nil {
		return handleEmbellishmentVariant(emb.Strike, emb.Thumb, emb.Light)
	}
	if staffSym.Pele != nil {
		return handleEmbellishment(emb.Pele)
	}
	if staffSym.LightPele != nil {
		return handleEmbellishmentVariant(emb.Pele, emb.NoVariant, emb.Light)
	}
	if staffSym.ThumbPele != nil {
		return handleEmbellishmentVariant(emb.Pele, emb.Thumb, emb.NoWeight)
	}
	if staffSym.LightThumbPele != nil {
		return handleEmbellishmentVariant(emb.Pele, emb.Thumb, emb.Light)
	}
	if staffSym.HalfPele != nil {
		return handleEmbellishmentVariant(emb.Pele, emb.Half, emb.NoWeight)
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
		return handleEmbellishmentVariant(emb.Pele, emb.Half, emb.Light)
	}
	if staffSym.GBirl != nil ||
		staffSym.ThumbBirl != nil {
		return handleEmbellishment(emb.GraceBirl)
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
		return handleMovement(movement.Embari, staffSym.Embari, true, true)
	}
	if staffSym.Endari != nil {
		return handleMovement(movement.Endari, staffSym.Endari, true, true)
	}
	if staffSym.Chedari != nil {
		return handleMovement(movement.Chedari, staffSym.Chedari, true, true)
	}
	if staffSym.Hedari != nil {
		return handleMovement(movement.Hedari, staffSym.Hedari, true, false)
	}
	if staffSym.Dili != nil {
		return handleMovement(movement.Dili, staffSym.Dili, true, true)
	}
	if staffSym.Tra != nil {
		return handleMovement(movement.Tra, staffSym.Tra, false, true)
	}
	if staffSym.Edre != nil {
		mv, _ := handleMovement(movement.Edre, staffSym.Edre, true, true)
		pitch := pitchFromSuffix(*staffSym.Edre)
		if pitch != common.E {
			mv.Note.Movement.PitchHint = pitch
		}
		return mv, nil
	}
	if staffSym.HalfEdre != nil {
		mv, _ := handleMovement(movement.Edre, staffSym.HalfEdre, true, true)
		mv.Note.Movement.Variant = movement.Half
		return mv, nil
	}
	if staffSym.GEdre != nil {
		return handleMovement(movement.Edre, staffSym.GEdre, true, true)
	}
	if staffSym.ThumbEdre != nil {
		return handleMovement(movement.Edre, staffSym.ThumbEdre, true, true)
	}
	if staffSym.Dare != nil {
		return handleMovement(movement.Dare, staffSym.Dare, true, true)
	}
	if staffSym.HalfDare != nil {
		return handleMovement(movement.Dare, staffSym.HalfDare, true, true)
	}
	if staffSym.ThumbDare != nil {
		return handleMovement(movement.Dare, staffSym.ThumbDare, true, true)
	}
	if staffSym.GDare != nil {
		return handleMovement(movement.Dare, staffSym.GDare, true, true)
	}
	if staffSym.CheCheRe != nil {
		return handleMovement(movement.CheCheRe, staffSym.CheCheRe, true, true)
	}
	if staffSym.HalfCheCheRe != nil {
		return handleMovement(movement.CheCheRe, staffSym.HalfCheCheRe, true, true)
	}
	if staffSym.ThumbCheCheRe != nil {
		return handleMovement(movement.CheCheRe, staffSym.ThumbCheCheRe, true, true)
	}
	if staffSym.GripAbbrev != nil {
		return handleMovement(movement.Grip, staffSym.GripAbbrev, true, true)
	}
	if staffSym.Deda != nil {
		return handleMovement(movement.Deda, staffSym.Deda, true, true)
	}
	if staffSym.Enbain != nil {
		return handleMovement(movement.Enbain, staffSym.Enbain, true, true)
	}
	if staffSym.GEnbain != nil {
		return handleMovement(movement.Enbain, staffSym.GEnbain, true, true)
	}
	if staffSym.ThumbEnbain != nil {
		return handleMovement(movement.Enbain, staffSym.ThumbEnbain, true, true)
	}
	if staffSym.Otro != nil {
		return handleMovement(movement.Otro, staffSym.Otro, true, true)
	}
	if staffSym.GOtro != nil {
		return handleMovement(movement.Otro, staffSym.GOtro, true, true)
	}
	if staffSym.ThumbOtro != nil {
		return handleMovement(movement.Otro, staffSym.ThumbOtro, true, true)
	}
	if staffSym.Odro != nil {
		return handleMovement(movement.Odro, staffSym.Odro, true, true)
	}
	if staffSym.GOdro != nil {
		return handleMovement(movement.Odro, staffSym.GOdro, true, true)
	}
	if staffSym.ThumbOdro != nil {
		return handleMovement(movement.Odro, staffSym.ThumbOdro, true, true)
	}
	if staffSym.Adeda != nil {
		return handleMovement(movement.Adeda, staffSym.Adeda, true, true)
	}
	if staffSym.GAdeda != nil {
		return handleMovement(movement.Adeda, staffSym.GAdeda, true, true)
	}
	if staffSym.ThumbAdeda != nil {
		return handleMovement(movement.Adeda, staffSym.ThumbAdeda, true, true)
	}
	if staffSym.EchoBeats != nil {
		mv, _ := handleMovement(movement.EchoBeat, staffSym.EchoBeats, false, false)
		pitch := pitchFromSuffix(*staffSym.EchoBeats)
		mv.Note.Movement.Pitch = pitch
		return mv, nil
	}
	if staffSym.Darodo != nil {
		return handleMovement(movement.Darodo, staffSym.Darodo, false, true)
	}
	if staffSym.Hiharin != nil {
		return handleMovement(movement.Hiharin, staffSym.Hiharin, false, false)
	}
	if staffSym.Rodin != nil {
		return handleMovement(movement.Rodin, staffSym.Rodin, false, false)
	}
	if staffSym.Chelalho != nil {
		return handleMovement(movement.Chelalho, staffSym.Chelalho, false, false)
	}
	if staffSym.Din != nil {
		return handleMovement(movement.Din, staffSym.Din, false, false)
	}
	if staffSym.Lemluath != nil {
		lemSym, hadBrea := stripBreabach(staffSym.Lemluath)
		mv, _ := handleMovement(movement.Lemluath, &lemSym, false, true)
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
		sym, _ := handleMovement(movement.Lemluath, &lemSym, false, true)
		move := sym.Note.Movement
		pitch := pitchFromSuffix(lemSym)
		move.PitchHint = pitch
		move.Breabach = hadBrea
		lastSym.Note.Movement = move
		return nil, nil
	}
	if staffSym.TaorluathPio != nil {
		return handleMovementWithPitchHintSuffixAndBreabach(
			movement.Taorluath, staffSym.TaorluathPio, false, true,
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
			movement.Taorluath, staffSym.TaorluathAbbrev, false, true,
		)
		move := sym.Note.Movement
		lastSym.Note.Movement = move
		return nil, err
	}
	if staffSym.TaorluathAmach != nil {
		pitch := pitchFromSuffix(*staffSym.TaorluathAmach)
		mv, _ := handleMovement(movement.Taorluath, staffSym.TaorluathAmach, false, false)
		mv.Note.Movement.AMach = true
		mv.Note.Movement.PitchHint = pitch
		return mv, nil
	}
	if staffSym.Crunluath != nil {
		mv, _ := handleMovementWithPitchHintSuffixAndBreabach(
			movement.Crunluath, staffSym.Crunluath, false, true,
		)
		if strings.Contains(*staffSym.Crunluath, "crunllgla") {
			mv.Note.Movement.AdditionalPitchHint = common.LowG
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
			movement.Crunluath, staffSym.CrunluathAbbrev, false, true,
		)
		move := sym.Note.Movement
		lastSym.Note.Movement = move
		return nil, err
	}
	if staffSym.CrunluathAmach != nil {
		pitch := pitchFromSuffix(*staffSym.CrunluathAmach)
		mv, _ := handleMovement(movement.Crunluath, staffSym.CrunluathAmach, false, false)
		mv.Note.Movement.AMach = true
		mv.Note.Movement.PitchHint = pitch
		return mv, nil
	}
	if staffSym.Tripling != nil {
		pitch := pitchFromSuffix(*staffSym.Tripling)
		mv, _ := handleMovement(movement.Tripling, staffSym.Tripling, false, true)
		// thumb variant handled here because pt is always recognized as thumb
		if strings.HasPrefix(*staffSym.Tripling, "ptt") {
			mv.Note.Movement.Variant = movement.Thumb
		}
		mv.Note.Movement.Pitch = pitch
		return mv, nil
	}
	if staffSym.Tempo != nil {
		tempo, err := strconv.ParseUint(staffSym.Tempo.Tempo, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed parsing tune tempo: %s", err.Error())
		}
		return &music_model.Symbol{TempoChange: tempo}, nil
	}

	return nil, nil // fmt.Errorf("staff symbol %v not handled", staffSym)
}

func handleMovementWithPitchHintSuffixAndBreabach(
	mtype movement.Type,
	sym *string,
	withThumb bool,
	withHalf bool,
) (*music_model.Symbol, error) {
	strippedSym, hadBrea := stripBreabach(sym)
	currSym, _ := handleMovement(mtype, &strippedSym, withThumb, withHalf)
	pitch := pitchFromSuffix(strippedSym)
	currSym.Note.Movement.PitchHint = pitch
	currSym.Note.Movement.Breabach = hadBrea

	return currSym, nil
}

func handleMovement(mtype movement.Type, sym *string, withThumb bool, withHalf bool) (*music_model.Symbol, error) {
	showAbbr := false
	if strings.HasPrefix(*sym, "p") {
		showAbbr = true
	}
	mVar := movement.NoVariant
	if withHalf {
		if strings.HasPrefix(*sym, "h") || strings.HasPrefix(*sym, "ph") {
			mVar = movement.Half
		}
	}

	if withThumb {
		if strings.HasPrefix(*sym, "t") || strings.HasPrefix(*sym, "pt") {
			mVar = movement.Thumb
		}
	}

	if strings.HasPrefix(*sym, "g") {
		mVar = movement.G
	}

	if strings.HasSuffix(*sym, "8") ||
		strings.HasSuffix(*sym, "16") {
		mVar = movement.LongLowG
	}

	return &music_model.Symbol{
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
) (*music_model.Symbol, error) {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Movement: &movement.Movement{
				Type:    movement.Cadence,
				Fermata: fermata,
				Pitches: pitchesFromCadenceSym(*cad, fermata),
			},
		},
	}, nil
}

func pitchesFromCadenceSym(sym string, fermata bool) []common.Pitch {
	if fermata {
		sym = strings.Replace(sym, "fcad", "", 1)
	} else {
		sym = strings.Replace(sym, "cad", "", 1)
	}

	pitches := make([]common.Pitch, len(sym))
	for i, ch := range sym {
		switch ch {
		case 'g':
			pitches[i] = common.HighG
		case 'e':
			pitches[i] = common.E
		case 'd':
			pitches[i] = common.D
		case 'a':
			pitches[i] = common.HighA
		case 'f':
			pitches[i] = common.F
		default:
			log.Error().Msgf("char %c is not handled for cadence symbol", ch)
			pitches[i] = common.NoPitch
		}
	}
	return pitches
}

func handleInsideStaffComment(
	lastSym *music_model.Symbol,
	currentMeasure *music_model.Measure,
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
	embType emb.EmbellishmentType,
) (*music_model.Symbol, error) {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Embellishment: &emb.Embellishment{
				Type: embType,
			},
		},
	}, nil
}

func handleDoubleGrace(pitch common.Pitch) (*music_model.Symbol, error) {
	doubleG, err := handleEmbellishment(emb.DoubleGrace)
	doubleG.Note.Embellishment.Pitch = pitch
	return doubleG, err
}

func handleEmbellishmentVariant(
	embType emb.EmbellishmentType,
	variant emb.EmbellishmentVariant,
	weight emb.EmbellishmentWeight,
) (*music_model.Symbol, error) {
	return &music_model.Symbol{
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

func handleAccidential(acc accidental.Accidental) *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Accidental: acc,
		},
	}
}

func handleTimeLine(sym string) (*music_model.Symbol, error) {
	if sym == "'1" {
		return newTimeLineStartSymbol(time_line.First), nil
	}
	if sym == "'2" {
		return newTimeLineStartSymbol(time_line.Second), nil
	}
	if sym == "'22" {
		return newTimeLineStartSymbol(time_line.SecondOf2), nil
	}
	if sym == "'23" {
		return newTimeLineStartSymbol(time_line.SecondOf3), nil
	}
	if sym == "'24" {
		return newTimeLineStartSymbol(time_line.SecondOf4), nil
	}
	if sym == "'224" {
		return newTimeLineStartSymbol(time_line.SecondOf2And4), nil
	}
	if sym == "'25" {
		return newTimeLineStartSymbol(time_line.SecondOf5), nil
	}
	if sym == "'26" {
		return newTimeLineStartSymbol(time_line.SecondOf6), nil
	}
	if sym == "'27" {
		return newTimeLineStartSymbol(time_line.SecondOf7), nil
	}
	if sym == "'28" {
		return newTimeLineStartSymbol(time_line.SecondOf8), nil
	}
	if sym == "'si" {
		return newTimeLineStartSymbol(time_line.Singling), nil
	}
	if sym == "'do" {
		return newTimeLineStartSymbol(time_line.Doubling), nil
	}
	if sym == "'bis" {
		return newTimeLineStartSymbol(time_line.Bis), nil
	}
	if sym == "'intro" {
		return newTimeLineStartSymbol(time_line.Intro), nil
	}

	return nil, fmt.Errorf("time line symbol %s not handled", sym)
}

func newTimeLineStartSymbol(ttype time_line.TimeLineType) *music_model.Symbol {
	return &music_model.Symbol{
		TimeLine: &time_line.TimeLine{
			BoundaryType: time_line.Start,
			Type:         ttype,
		},
	}
}

func embellishmentForSingleGrace(grace *string) *emb.Embellishment {
	emb := &emb.Embellishment{
		Type: emb.SingleGrace,
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

func pitchFromSuffix(sym string) common.Pitch {
	if strings.HasSuffix(sym, "lg") {
		return common.LowG
	}
	if strings.HasSuffix(sym, "la") {
		return common.LowA
	}
	if strings.HasSuffix(sym, "b") {
		return common.B
	}
	if strings.HasSuffix(sym, "c") {
		return common.C
	}
	if strings.HasSuffix(sym, "d") {
		return common.D
	}
	if strings.HasSuffix(sym, "e") {
		return common.E
	}
	if strings.HasSuffix(sym, "f") {
		return common.F
	}
	if strings.HasSuffix(sym, "hg") {
		return common.HighG
	}
	if strings.HasSuffix(sym, "ha") {
		return common.HighA
	}
	return common.NoPitch
}

func stripBreabach(sym *string) (string, bool) {
	stripped := strings.Replace(*sym, "brea", "", 1)
	didReplace := len(*sym) != len(stripped)
	return stripped, didReplace
}

func newTimeLineEnd(sym *string) *music_model.Symbol {
	ttype := time_line.NoType
	if *sym == "bis_'" {
		ttype = time_line.Bis
	}
	return &music_model.Symbol{
		TimeLine: &time_line.TimeLine{
			BoundaryType: time_line.End,
			Type:         ttype,
		},
	}
}
