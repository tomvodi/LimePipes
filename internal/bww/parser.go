package bww

import (
	"banduslib/internal/common"
	"banduslib/internal/interfaces"
	"github.com/alecthomas/participle/v2"
)

type bwwParser struct {
}

func (b *bwwParser) ParseBwwData(data []byte) (*common.BwwDocument, error) {
	parser, err := participle.Build[common.BwwDocument](
		participle.Elide("WHITESPACE"),
		participle.Lexer(BwwLexer),
		participle.Unquote("STRING"),
	)
	if err != nil {
		return nil, err
	}

	var bwwDoc *common.BwwDocument
	bwwDoc, err = parser.ParseBytes("", data)
	if err != nil {
		return nil, err
	}

	return bwwDoc, nil
}

func NewBwwParser() interfaces.BwwParser {
	return &bwwParser{}
}
