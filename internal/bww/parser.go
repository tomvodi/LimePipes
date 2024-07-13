package bww

import (
	"github.com/alecthomas/participle/v2"
	"github.com/tomvodi/limepipes/internal/common/music_model"
	"github.com/tomvodi/limepipes/internal/interfaces"
)

type bwwParser struct {
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

	return convertGrammarToModel(bwwDoc)
}

func NewBwwParser() interfaces.BwwParser {
	return &bwwParser{}
}
