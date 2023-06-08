package expander

import (
	c "banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/symbols/embellishment"
	"banduslib/internal/interfaces"
)

type birlsExp struct {
}

func (b *birlsExp) ExpandSymbol(symbol *music_model.Symbol, prevSymPitch c.Pitch) {
	if symbol == nil || symbol.Note == nil || symbol.Note.Embellishment == nil {
		return
	}

	emb := symbol.Note.Embellishment

	var expanded []c.Pitch
	isHalf := prevSymPitch == c.LowA
	if isHalf {
		expanded = []c.Pitch{c.LowG, c.LowA, c.LowG}
	}
	if emb.Variant == embellishment.G {
		expanded = []c.Pitch{c.HighG, c.LowA, c.LowG, c.LowA, c.LowG}
	}
	if emb.Variant == embellishment.Thumb {
		expanded = []c.Pitch{c.HighA, c.LowA, c.LowG, c.LowA, c.LowG}
	}
	if expanded == nil {
		expanded = []c.Pitch{c.LowA, c.LowG, c.LowA, c.LowG}
	}
	symbol.Note.ExpandedEmbellishment = expanded
}

func NewBirlsExpander() interfaces.SymbolExpander {
	return &birlsExp{}
}
