package expander

import (
	c "banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/symbols/embellishment"
	"banduslib/internal/interfaces"
)

type dblStrikeExp struct {
}

func (d *dblStrikeExp) ExpandSymbol(symbol *music_model.Symbol, _ c.Pitch) {
	if symbol == nil || symbol.Note == nil || symbol.Note.Embellishment == nil {
		return
	}

	emb := symbol.Note.Embellishment

	isLight := emb.Weight == embellishment.Light
	var basicDblStrike []c.Pitch
	var symbolPitch = symbol.Note.Pitch
	if symbolPitch >= c.LowA &&
		symbolPitch <= c.D {
		basicDblStrike = []c.Pitch{c.LowG, symbolPitch, c.LowG}
		if isLight {
			basicDblStrike[0] = c.C
			basicDblStrike[2] = c.C
		}
	}
	if symbolPitch == c.E {
		basicDblStrike = []c.Pitch{c.LowA, symbolPitch, c.LowA}
	}
	if symbolPitch == c.F {
		basicDblStrike = []c.Pitch{c.E, symbolPitch, c.E}
	}
	if symbolPitch == c.HighG {
		basicDblStrike = []c.Pitch{c.F, symbolPitch, c.F}
	}
	if symbolPitch == c.HighA &&
		(emb.Variant == embellishment.Half || emb.Variant == embellishment.NoVariant) {
		basicDblStrike = []c.Pitch{c.HighG, symbolPitch, c.HighG}
	}
	if emb.Variant != embellishment.NoVariant {
		basicDblStrike = append([]c.Pitch{symbolPitch}, basicDblStrike...)
	}

	if emb.Variant == embellishment.Thumb {
		basicDblStrike = append([]c.Pitch{c.HighA}, basicDblStrike...)
	}
	if emb.Variant == embellishment.G {
		basicDblStrike = append([]c.Pitch{c.HighG}, basicDblStrike...)
	}

	symbol.Note.ExpandedEmbellishment = basicDblStrike
}

func NewDoubleStrikesExpander() interfaces.SymbolExpander {
	return &dblStrikeExp{}
}
