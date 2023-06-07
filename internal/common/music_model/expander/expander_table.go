package expander

import (
	emb "banduslib/internal/common/music_model/symbols/embellishment"
	"banduslib/internal/interfaces"
)

type ExpandTable map[emb.Embellishment]interfaces.SymbolExpander

func newSymbolExpanderTable() ExpandTable {
	dblExpander := NewDoublingsExpander()
	return map[emb.Embellishment]interfaces.SymbolExpander{
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
