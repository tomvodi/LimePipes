package expander

import (
	"banduslib/internal/common"
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
	prevSymPitch := common.NoPitch
	for _, measure := range tune.Measures {
		for _, symbol := range measure.Symbols {
			e.expandSymbol(symbol, prevSymPitch)

			if symbol.IsValidNote() {
				prevSymPitch = symbol.Note.Pitch
			} else {
				prevSymPitch = common.NoPitch
			}
		}
	}
}

func (e *embExpander) expandSymbol(symbol *music_model.Symbol, prevSymPitch common.Pitch) {
	if !symbol.IsValidNote() {
		return
	}

	if symbol.Note.Embellishment == nil {
		return
	}

	expander, ok := e.table[*symbol.Note.Embellishment]
	if !ok {
		log.Error().Msgf("no embellishment expander for %v", *symbol.Note.Embellishment)
		return
	}

	expander.ExpandSymbol(symbol, prevSymPitch)
}

func NewEmbellishmentExpander() interfaces.EmbellishmentExpander {
	table := newSymbolExpanderTable()
	return &embExpander{
		table: table,
	}
}
