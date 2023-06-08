package expander

import (
	"banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/interfaces"
)

type singleGraceExpander struct {
}

func (s *singleGraceExpander) ExpandSymbol(symbol *music_model.Symbol, _ common.Pitch) {
	if symbol == nil || symbol.Note == nil || symbol.Note.Embellishment == nil {
		return
	}

	emb := symbol.Note.Embellishment
	symbol.Note.ExpandedEmbellishment = []common.Pitch{
		emb.Pitch,
	}
}

func NewSingleGraceExpander() interfaces.SymbolExpander {
	return &singleGraceExpander{}
}
