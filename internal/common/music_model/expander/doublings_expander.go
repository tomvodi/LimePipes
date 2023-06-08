package expander

import (
	"banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/symbols"
	"banduslib/internal/common/music_model/symbols/embellishment"
	"banduslib/internal/interfaces"
)

type dblExpand struct {
}

func (d *dblExpand) ExpandSymbol(symbol *music_model.Symbol, _ common.Pitch) {
	if symbol == nil || symbol.Note == nil || symbol.Note.Embellishment == nil {
		return
	}

	emb := symbol.Note.Embellishment

	switch emb.Variant {
	case embellishment.NoVariant:
		handleRegular(symbol.Note)
	case embellishment.Thumb:
		handleThumb(symbol.Note)
	case embellishment.Half:
		handleHalf(symbol.Note)
	}
}

func handleRegular(note *symbols.Note) {
	if note.Pitch >= common.HighG {
		note.ExpandedEmbellishment = []common.Pitch{
			note.Pitch,
			note.Pitch - 1,
		}
		return
	}

	if note.Pitch >= common.D {
		note.ExpandedEmbellishment = []common.Pitch{
			common.HighG,
			note.Pitch,
			note.Pitch + 1,
		}
		return
	}

	note.ExpandedEmbellishment = []common.Pitch{
		common.HighG,
		note.Pitch,
		common.D,
	}
}

func handleThumb(note *symbols.Note) {
	if note.Pitch >= common.HighG {
		return
	}

	if note.Pitch >= common.D {
		note.ExpandedEmbellishment = []common.Pitch{
			common.HighA,
			note.Pitch,
			note.Pitch + 1,
		}
		return
	}

	note.ExpandedEmbellishment = []common.Pitch{
		common.HighA,
		note.Pitch,
		common.D,
	}
}

func handleHalf(note *symbols.Note) {
	if note.Pitch >= common.HighG {
		return
	}

	if note.Pitch >= common.D {
		note.ExpandedEmbellishment = []common.Pitch{
			note.Pitch,
			note.Pitch + 1,
		}
		return
	}

	note.ExpandedEmbellishment = []common.Pitch{
		note.Pitch,
		common.D,
	}
}

func NewDoublingsExpander() interfaces.SymbolExpander {
	return &dblExpand{}
}
