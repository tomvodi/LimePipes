package expander

import (
	c "banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/interfaces"
)

type taorExpand struct {
}

func (t *taorExpand) ExpandSymbol(symbol *music_model.Symbol, prevSymPitch c.Pitch) {
	if symbol == nil || symbol.Note == nil || symbol.Note.Embellishment == nil {
		return
	}

	emb := symbol.Note.Embellishment

	// regular taorluath
	var expanded []c.Pitch
	isHalf := prevSymPitch == c.LowG
	isB := emb.Pitch == c.B
	if isHalf {
		expanded = []c.Pitch{c.D, c.LowG, c.E}
	}
	if isB {
		expanded = []c.Pitch{c.LowG, c.B, c.LowG, c.E}
	}
	if expanded == nil {
		expanded = []c.Pitch{c.LowG, c.D, c.LowG, c.E}
	}
	symbol.Note.ExpandedEmbellishment = expanded

}

func NewTaorluathsExpander() interfaces.SymbolExpander {
	return &taorExpand{}
}
