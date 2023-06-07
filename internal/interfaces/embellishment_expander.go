package interfaces

import "banduslib/internal/common/music_model"

//go:generate mockgen -source embellishment_expander.go -destination ./mocks/mock_embellishment_expander.go.go

type SymbolExpander interface {
	// ExpandSymbol expands all embellishments in the music model symbol
	ExpandSymbol(symbol *music_model.Symbol)
}

type EmbellishmentExpander interface {
	SymbolExpander
	// ExpandModel expands all embellishments in music model
	ExpandModel(model music_model.MusicModel)

	// ExpandTune expands all embellishments in music model tune
	ExpandTune(tune *music_model.Tune)
}
