package expander

import (
	c "banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/symbols/embellishment"
	"banduslib/internal/interfaces"
)

type throwdExpand struct {
}

func (t *throwdExpand) ExpandSymbol(symbol *music_model.Symbol, prevSymPitch c.Pitch) {
	if symbol == nil || symbol.Note == nil || symbol.Note.Embellishment == nil {
		return
	}

	emb := symbol.Note.Embellishment

	var expanded []c.Pitch
	if emb.Weight == embellishment.Light {
		expanded = []c.Pitch{c.LowG, c.D, c.C}
	}
	if emb.Weight == embellishment.Heavy {
		expanded = []c.Pitch{c.LowG, c.D, c.LowG, c.C}
	}

	isHalf := prevSymPitch == c.LowG
	if isHalf {
		expanded = expanded[1:]
	}

	symbol.Note.ExpandedEmbellishment = expanded

}

func NewThrowdsExpander() interfaces.SymbolExpander {
	return &throwdExpand{}
}
