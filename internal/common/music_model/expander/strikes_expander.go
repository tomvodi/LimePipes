package expander

import (
	"banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/symbols/embellishment"
	"banduslib/internal/interfaces"
)

type strikesExp struct {
}

func (s *strikesExp) ExpandSymbol(symbol *music_model.Symbol, _ common.Pitch) {
	if symbol == nil || symbol.Note == nil || symbol.Note.Embellishment == nil {
		return
	}

	emb := symbol.Note.Embellishment
	if emb.Variant == embellishment.NoVariant {
		// Single strikes
		symbol.Note.ExpandedEmbellishment = []common.Pitch{emb.Pitch}
	}
	if emb.Variant == embellishment.Half {
		symbol.Note.ExpandedEmbellishment =
			halfStrikePitchesForWeight(symbol.Note.Pitch, emb.Weight)
	}
	if emb.Variant == embellishment.Thumb {
		pitchesHalf := halfStrikePitchesForWeight(symbol.Note.Pitch, emb.Weight)
		pitchesHalf = append([]common.Pitch{common.HighA}, pitchesHalf...)
		symbol.Note.ExpandedEmbellishment = pitchesHalf
	}
	if emb.Variant == embellishment.G {
		pitchesHalf := halfStrikePitchesForWeight(symbol.Note.Pitch, emb.Weight)
		pitchesHalf = append([]common.Pitch{common.HighG}, pitchesHalf...)
		symbol.Note.ExpandedEmbellishment = pitchesHalf
	}
}

func halfStrikePitchesForWeight(
	pitch common.Pitch,
	weight embellishment.EmbellishmentWeight,
) []common.Pitch {
	pitches := []common.Pitch{pitch}
	if pitch >= common.LowA && pitch <= common.C {
		pitches = append(pitches, common.LowG)
	}

	if pitch == common.D {
		if weight == embellishment.Light {
			pitches = append(pitches, common.C)
		} else {
			pitches = append(pitches, common.LowG)
		}
	}
	if pitch == common.E {
		pitches = append(pitches, common.LowA)
	}
	if pitch >= common.F && pitch <= common.HighG {
		pitches = append(pitches, pitch-1)
	}

	return pitches
}

func NewStrikesExpander() interfaces.SymbolExpander {
	return &strikesExp{}
}
