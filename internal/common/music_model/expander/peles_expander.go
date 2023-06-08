package expander

import (
	c "banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/symbols/embellishment"
	"banduslib/internal/interfaces"
)

type peleExp struct {
}

func (p *peleExp) ExpandSymbol(symbol *music_model.Symbol, prevSymPitch c.Pitch) {
	if symbol == nil || symbol.Note == nil || symbol.Note.Embellishment == nil {
		return
	}

	emb := symbol.Note.Embellishment

	isLight := emb.Weight == embellishment.Light
	var basicPele []c.Pitch
	var symbolPitch = symbol.Note.Pitch
	if symbolPitch >= c.LowA &&
		symbolPitch <= c.D {
		basicPele = []c.Pitch{symbolPitch, c.E, symbolPitch, c.LowG}
		if isLight {
			basicPele[3] = c.C
		}
	}
	if symbolPitch == c.E {
		basicPele = []c.Pitch{symbolPitch, c.F, symbolPitch, c.LowA}
	}
	if symbolPitch == c.F {
		basicPele = []c.Pitch{symbolPitch, c.HighG, symbolPitch, c.E}
	}
	if symbolPitch == c.HighG &&
		(emb.Variant == embellishment.Thumb || emb.Variant == embellishment.Half) {
		basicPele = []c.Pitch{symbolPitch, c.HighA, symbolPitch, c.F}
	}

	if emb.Variant == embellishment.Thumb {
		basicPele = append([]c.Pitch{c.HighA}, basicPele...)
	}

	if emb.Variant == embellishment.NoVariant {
		basicPele = append([]c.Pitch{c.HighG}, basicPele...)
	}

	symbol.Note.ExpandedEmbellishment = basicPele
}

func NewPelesExpander() interfaces.SymbolExpander {
	return &peleExp{}
}
