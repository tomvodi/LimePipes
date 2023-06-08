package expander

import (
	c "banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/symbols/embellishment"
	"banduslib/internal/interfaces"
)

type trplStrikeExp struct {
}

func (t *trplStrikeExp) ExpandSymbol(symbol *music_model.Symbol, _ c.Pitch) {
	if symbol == nil || symbol.Note == nil || symbol.Note.Embellishment == nil {
		return
	}

	emb := symbol.Note.Embellishment

	isLight := emb.Weight == embellishment.Light
	var basicTrplStrike []c.Pitch
	var symbolPitch = symbol.Note.Pitch
	if symbolPitch >= c.LowA &&
		symbolPitch <= c.D {
		basicTrplStrike = []c.Pitch{c.LowG, symbolPitch, c.LowG, symbolPitch, c.LowG}
		if isLight {
			basicTrplStrike[0] = c.C
			basicTrplStrike[2] = c.C
			basicTrplStrike[4] = c.C
		}
	}
	if symbolPitch == c.E {
		basicTrplStrike = []c.Pitch{c.LowA, symbolPitch, c.LowA, symbolPitch, c.LowA}
	}
	if symbolPitch == c.F {
		basicTrplStrike = []c.Pitch{c.E, symbolPitch, c.E, symbolPitch, c.E}
	}
	if symbolPitch == c.HighG {
		basicTrplStrike = []c.Pitch{c.F, symbolPitch, c.F, symbolPitch, c.F}
	}
	if symbolPitch == c.HighA &&
		(emb.Variant == embellishment.Half || emb.Variant == embellishment.NoVariant) {
		basicTrplStrike = []c.Pitch{c.HighG, symbolPitch, c.HighG, symbolPitch, c.HighG}
	}
	if emb.Variant != embellishment.NoVariant {
		basicTrplStrike = append([]c.Pitch{symbolPitch}, basicTrplStrike...)
	}

	if emb.Variant == embellishment.Thumb {
		basicTrplStrike = append([]c.Pitch{c.HighA}, basicTrplStrike...)
	}
	if emb.Variant == embellishment.G {
		basicTrplStrike = append([]c.Pitch{c.HighG}, basicTrplStrike...)
	}

	symbol.Note.ExpandedEmbellishment = basicTrplStrike
}

func NewTripleStrikesExpander() interfaces.SymbolExpander {
	return &trplStrikeExp{}
}
