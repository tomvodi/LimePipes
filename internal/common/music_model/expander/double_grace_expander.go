package expander

import (
	"banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/interfaces"
	"github.com/rs/zerolog/log"
)

type dblGraceExpand struct {
}

func (d *dblGraceExpand) ExpandSymbol(symbol *music_model.Symbol, _ common.Pitch) {
	if symbol == nil || symbol.Note == nil || symbol.Note.Embellishment == nil {
		return
	}

	if symbol.Note.Pitch == common.LowG {
		log.Error().Msg("can't play double grace on LowG")
		return
	}

	emb := symbol.Note.Embellishment

	if symbol.Note.Pitch > emb.Pitch {
		log.Error().Msgf("can't play double grace %s on a melody note with pitch %s",
			emb.Pitch.String(), symbol.Note.Pitch.String())
		return
	}

	symbol.Note.ExpandedEmbellishment = []common.Pitch{
		emb.Pitch,
		symbol.Note.Pitch - 1,
	}
}

func NewDoubleGraceExpander() interfaces.SymbolExpander {
	return &dblGraceExpand{}
}
