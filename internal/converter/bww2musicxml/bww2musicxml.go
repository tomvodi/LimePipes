package bww2musicxml

import (
	"banduslib/internal/bww"
	"banduslib/internal/interfaces"
	"banduslib/internal/musicxml/model"
)

type bww2xml struct {
}

func (b *bww2xml) Convert(bww *bww.BwwDocument) (*model.Score, error) {
	//TODO implement me
	panic("implement me")
}

func NewBww2MusicxmlConverter() interfaces.Bww2Musicxml {
	return &bww2xml{}
}
