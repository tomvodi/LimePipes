package expander

import (
	"banduslib/internal/common"
	emb "banduslib/internal/common/music_model/symbols/embellishment"
	"banduslib/internal/interfaces"
)

type ExpandTable map[emb.Embellishment]interfaces.SymbolExpander

func newSymbolExpanderTable() ExpandTable {
	singleGraceExp := NewSingleGraceExpander()
	dblExpander := NewDoublingsExpander()
	return map[emb.Embellishment]interfaces.SymbolExpander{
		emb.Embellishment{
			Type:  emb.SingleGrace,
			Pitch: common.LowA,
		}: singleGraceExp,
		emb.Embellishment{
			Type:  emb.SingleGrace,
			Pitch: common.B,
		}: singleGraceExp,
		emb.Embellishment{
			Type:  emb.SingleGrace,
			Pitch: common.C,
		}: singleGraceExp,
		emb.Embellishment{
			Type:  emb.SingleGrace,
			Pitch: common.D,
		}: singleGraceExp,
		emb.Embellishment{
			Type:  emb.SingleGrace,
			Pitch: common.E,
		}: singleGraceExp,
		emb.Embellishment{
			Type:  emb.SingleGrace,
			Pitch: common.F,
		}: singleGraceExp,
		emb.Embellishment{
			Type:  emb.SingleGrace,
			Pitch: common.HighG,
		}: singleGraceExp,
		emb.Embellishment{
			Type:  emb.SingleGrace,
			Pitch: common.HighA,
		}: singleGraceExp,
		emb.Embellishment{
			Type: emb.Doubling,
		}: dblExpander,
		emb.Embellishment{
			Type:    emb.Doubling,
			Variant: emb.Thumb,
		}: dblExpander,
		emb.Embellishment{
			Type:    emb.Doubling,
			Variant: emb.Half,
		}: dblExpander,
	}
}
