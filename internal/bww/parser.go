package bww

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/interfaces"
	"github.com/alecthomas/participle/v2"
)

type bwwParser struct {
}

func (b *bwwParser) ParseBwwData(data []byte) ([]*music_model.Tune, error) {
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
