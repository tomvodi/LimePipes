package expander

import (
	c "banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/symbols"
	emb "banduslib/internal/common/music_model/symbols/embellishment"
	"banduslib/internal/utils"
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
)

func regThrowd(weight emb.EmbellishmentWeight) *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Pitch:  c.C,
			Length: c.Quarter,
			Embellishment: &emb.Embellishment{
				Type:   emb.ThrowD,
				Weight: weight,
			},
		},
	}
}

func Test_throwdExpander_regular_ExpandSymbol(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		symbol    *music_model.Symbol
		prevPitch c.Pitch
		want      []c.Pitch
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: "regular light",
			prepare: func(f *fields) {
				f.symbol = regThrowd(emb.Light)
				f.want = []c.Pitch{c.LowG, c.D, c.C}
			},
		},
		{
			name: "regular light with previous low g => half throwd",
			prepare: func(f *fields) {
				f.symbol = regThrowd(emb.Light)
				f.prevPitch = c.LowG
				f.want = []c.Pitch{c.D, c.C}
			},
		},
		{
			name: "regular heavy",
			prepare: func(f *fields) {
				f.symbol = regThrowd(emb.Heavy)
				f.want = []c.Pitch{c.LowG, c.D, c.LowG, c.C}
			},
		},
		{
			name: "regular heavy with previous low g => half throwd",
			prepare: func(f *fields) {
				f.symbol = regThrowd(emb.Heavy)
				f.prevPitch = c.LowG
				f.want = []c.Pitch{c.D, c.LowG, c.C}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewThrowdsExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}
