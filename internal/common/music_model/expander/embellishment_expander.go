package expander

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/interfaces"
	"github.com/rs/zerolog/log"
)

type embExpander struct {
	table ExpandTable
}

func (e *embExpander) ExpandModel(model music_model.MusicModel) {
	for _, tune := range model {
		e.ExpandTune(tune)
	}
}

func (e *embExpander) ExpandTune(tune *music_model.Tune) {
	for _, measure := range tune.Measures {
		for _, symbol := range measure.Symbols {
			e.ExpandSymbol(symbol)
		}
	}
}

func (e *embExpander) ExpandSymbol(symbol *music_model.Symbol) {
	if !symbol.IsValidNote() {
		return
	}

	if symbol.Note.Embellishment == nil {
		return
	}

	packer, ok := e.table[*symbol.Note.Embellishment]
	if !ok {
		log.Error().Msgf("no embellishment expander for %v", *symbol.Note.Embellishment)
		return
	}

	packer.ExpandSymbol(symbol)
}

func NewEmbellishmentExpander() interfaces.EmbellishmentExpander {
	table := newSymbolExpanderTable()
	return &embExpander{
		table: table,
	}
}
