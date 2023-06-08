package expander

import (
	c "banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/interfaces"
)

type bubblysExpand struct {
}

func (b *bubblysExpand) ExpandSymbol(symbol *music_model.Symbol, prevSymPitch c.Pitch) {
	if symbol == nil || symbol.Note == nil || symbol.Note.Embellishment == nil {
		return
	}

	var expanded []c.Pitch
	isHalf := prevSymPitch == c.LowG
	if isHalf {
		expanded = []c.Pitch{c.D, c.LowG, c.C, c.LowG}
	} else {
		expanded = []c.Pitch{c.LowG, c.D, c.LowG, c.C, c.LowG}
	}
	symbol.Note.ExpandedEmbellishment = expanded
}

func NewBubblysExpander() interfaces.SymbolExpander {
	return &bubblysExpand{}
}
