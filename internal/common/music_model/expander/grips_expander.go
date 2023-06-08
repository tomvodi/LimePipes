package expander

import (
	c "banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/symbols/embellishment"
	"banduslib/internal/interfaces"
)

type grpExpander struct {
}

func (g *grpExpander) ExpandSymbol(symbol *music_model.Symbol, prevSymPitch c.Pitch) {
	if symbol == nil || symbol.Note == nil || symbol.Note.Embellishment == nil {
		return
	}

	emb := symbol.Note.Embellishment
	if emb.Variant == embellishment.NoVariant {
		// regular grip
		var expanded []c.Pitch
		isHalf := prevSymPitch == c.LowG
		isB := emb.Pitch == c.B
		if isHalf {
			expanded = []c.Pitch{c.D, c.LowG}
		}
		if isB {
			expanded = []c.Pitch{c.LowG, c.B, c.LowG}
		}
		if expanded == nil {
			expanded = []c.Pitch{c.LowG, c.D, c.LowG}
		}
		symbol.Note.ExpandedEmbellishment = expanded
		return
	}
	basicGrp := []c.Pitch{symbol.Note.Pitch, c.LowG, c.D, c.LowG}
	if emb.Pitch == c.B {
		basicGrp[2] = c.B
	}
	if symbol.Note.Pitch == c.F {
		basicGrp[2] = c.F
	}
	if symbol.Note.Pitch == c.HighG && emb.Variant == embellishment.Thumb {
		basicGrp[2] = c.F
	}

	if emb.Variant == embellishment.G {
		basicGrp = append([]c.Pitch{c.HighG}, basicGrp...)
	}
	if emb.Variant == embellishment.Thumb {
		basicGrp = append([]c.Pitch{c.HighA}, basicGrp...)
	}

	symbol.Note.ExpandedEmbellishment = basicGrp
}

func NewGripsExpander() interfaces.SymbolExpander {
	return &grpExpander{}
}
