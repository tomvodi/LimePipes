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
	strikesExpander := NewStrikesExpander()
	gripsExpander := NewGripsExpander()
	taorExpander := NewTaorluathsExpander()
	birlsExpander := NewBirlsExpander()
	throwdExpander := NewThrowdsExpander()
	pelesExpander := NewPelesExpander()
	doubleStrikesExpander := NewDoubleStrikesExpander()
	tripleStrikesExpander := NewTripleStrikesExpander()

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
		emb.Embellishment{
			Type:  emb.Strike,
			Pitch: common.LowG,
		}: strikesExpander,
		emb.Embellishment{
			Type:  emb.Strike,
			Pitch: common.LowA,
		}: strikesExpander,
		emb.Embellishment{
			Type:  emb.Strike,
			Pitch: common.B,
		}: strikesExpander,
		emb.Embellishment{
			Type:  emb.Strike,
			Pitch: common.C,
		}: strikesExpander,
		emb.Embellishment{
			Type:  emb.Strike,
			Pitch: common.D,
		}: strikesExpander,
		emb.Embellishment{
			Type:  emb.Strike,
			Pitch: common.E,
		}: strikesExpander,
		emb.Embellishment{
			Type:  emb.Strike,
			Pitch: common.F,
		}: strikesExpander,
		emb.Embellishment{
			Type:  emb.Strike,
			Pitch: common.HighG,
		}: strikesExpander,
		emb.Embellishment{
			Type:    emb.Strike,
			Variant: emb.G,
		}: strikesExpander,
		emb.Embellishment{
			Type:    emb.Strike,
			Variant: emb.G,
			Weight:  emb.Light,
		}: strikesExpander,
		emb.Embellishment{
			Type:    emb.Strike,
			Variant: emb.Thumb,
		}: strikesExpander,
		emb.Embellishment{
			Type:    emb.Strike,
			Variant: emb.Thumb,
			Weight:  emb.Light,
		}: strikesExpander,
		emb.Embellishment{
			Type:    emb.Strike,
			Variant: emb.Half,
		}: strikesExpander,
		emb.Embellishment{
			Type: emb.Grip,
		}: gripsExpander,
		emb.Embellishment{
			Type:  emb.Grip,
			Pitch: common.B,
		}: gripsExpander,
		emb.Embellishment{
			Type:    emb.Grip,
			Variant: emb.G,
		}: gripsExpander,
		emb.Embellishment{
			Type:    emb.Grip,
			Variant: emb.G,
			Pitch:   common.B,
		}: gripsExpander,
		emb.Embellishment{
			Type:    emb.Grip,
			Variant: emb.Thumb,
		}: gripsExpander,
		emb.Embellishment{
			Type:    emb.Grip,
			Variant: emb.Thumb,
			Pitch:   common.B,
		}: gripsExpander,
		emb.Embellishment{
			Type:    emb.Grip,
			Variant: emb.Half,
		}: gripsExpander,
		emb.Embellishment{
			Type:    emb.Grip,
			Variant: emb.Half,
			Pitch:   common.B,
		}: gripsExpander,
		emb.Embellishment{
			Type: emb.Taorluath,
		}: taorExpander,
		emb.Embellishment{
			Type:  emb.Taorluath,
			Pitch: common.B,
		}: taorExpander,
		emb.Embellishment{
			Type: emb.Bubbly,
		}: NewBubblysExpander(),
		emb.Embellishment{
			Type: emb.Birl,
		}: birlsExpander,
		emb.Embellishment{
			Type:    emb.Birl,
			Variant: emb.G,
		}: birlsExpander,
		emb.Embellishment{
			Type:    emb.Birl,
			Variant: emb.Thumb,
		}: birlsExpander,
		emb.Embellishment{
			Type:   emb.ThrowD,
			Weight: emb.Light,
		}: throwdExpander,
		emb.Embellishment{
			Type:   emb.ThrowD,
			Weight: emb.Heavy,
		}: throwdExpander,
		emb.Embellishment{
			Type: emb.Pele,
		}: pelesExpander,
		emb.Embellishment{
			Type:   emb.Pele,
			Weight: emb.Light,
		}: pelesExpander,
		emb.Embellishment{
			Type:    emb.Pele,
			Variant: emb.Thumb,
		}: pelesExpander,
		emb.Embellishment{
			Type:    emb.Pele,
			Variant: emb.Thumb,
			Weight:  emb.Light,
		}: pelesExpander,
		emb.Embellishment{
			Type:    emb.Pele,
			Variant: emb.Half,
		}: pelesExpander,
		emb.Embellishment{
			Type:    emb.Pele,
			Variant: emb.Half,
			Weight:  emb.Light,
		}: pelesExpander,
		emb.Embellishment{
			Type: emb.DoubleStrike,
		}: doubleStrikesExpander,
		emb.Embellishment{
			Type:   emb.DoubleStrike,
			Weight: emb.Light,
		}: doubleStrikesExpander,
		emb.Embellishment{
			Type:    emb.DoubleStrike,
			Variant: emb.G,
		}: doubleStrikesExpander,
		emb.Embellishment{
			Type:    emb.DoubleStrike,
			Variant: emb.G,
			Weight:  emb.Light,
		}: doubleStrikesExpander,
		emb.Embellishment{
			Type:    emb.DoubleStrike,
			Variant: emb.Thumb,
		}: doubleStrikesExpander,
		emb.Embellishment{
			Type:    emb.DoubleStrike,
			Variant: emb.Thumb,
			Weight:  emb.Light,
		}: doubleStrikesExpander,
		emb.Embellishment{
			Type:    emb.DoubleStrike,
			Variant: emb.Half,
		}: doubleStrikesExpander,
		emb.Embellishment{
			Type:    emb.DoubleStrike,
			Variant: emb.Half,
			Weight:  emb.Light,
		}: doubleStrikesExpander,
		emb.Embellishment{
			Type: emb.TripleStrike,
		}: tripleStrikesExpander,
		emb.Embellishment{
			Type:   emb.TripleStrike,
			Weight: emb.Light,
		}: tripleStrikesExpander,
		emb.Embellishment{
			Type:    emb.TripleStrike,
			Variant: emb.G,
		}: tripleStrikesExpander,
		emb.Embellishment{
			Type:    emb.TripleStrike,
			Variant: emb.G,
			Weight:  emb.Light,
		}: tripleStrikesExpander,
		emb.Embellishment{
			Type:    emb.TripleStrike,
			Variant: emb.Thumb,
		}: tripleStrikesExpander,
		emb.Embellishment{
			Type:    emb.TripleStrike,
			Variant: emb.Thumb,
			Weight:  emb.Light,
		}: tripleStrikesExpander,
		emb.Embellishment{
			Type:    emb.TripleStrike,
			Variant: emb.Half,
		}: tripleStrikesExpander,
		emb.Embellishment{
			Type:    emb.TripleStrike,
			Variant: emb.Half,
			Weight:  emb.Light,
		}: tripleStrikesExpander,
	}
}
