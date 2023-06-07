package bww

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/interfaces"
	"github.com/alecthomas/participle/v2"
)

type bwwParser struct {
	embExpander interfaces.EmbellishmentExpander
}

func (b *bwwParser) ParseBwwData(data []byte) (music_model.MusicModel, error) {
	parser, err := participle.Build[BwwDocument](
		participle.Elide("WHITESPACE"),
		participle.Lexer(BwwLexer),
		participle.Unquote("STRING"),
	)
	if err != nil {
		return nil, err
	}

	var bwwDoc *BwwDocument
	bwwDoc, err = parser.ParseBytes("", data)
	if err != nil {
		return nil, err
	}

	muMo, err := convertGrammarToModel(bwwDoc)
	if err != nil {
		return nil, err
	}

	b.embExpander.ExpandModel(muMo)
	return muMo, nil
}

func NewBwwParser(
	embExpander interfaces.EmbellishmentExpander,
) interfaces.BwwParser {
	return &bwwParser{
		embExpander: embExpander,
	}
}
