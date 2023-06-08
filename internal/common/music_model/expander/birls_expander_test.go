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

func regBirl() *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Pitch:  c.LowA,
			Length: c.Quarter,
			Embellishment: &emb.Embellishment{
				Type: emb.Birl,
			},
		},
	}
}

func gBirl() *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Pitch:  c.LowA,
			Length: c.Quarter,
			Embellishment: &emb.Embellishment{
				Type:    emb.Birl,
				Variant: emb.G,
			},
		},
	}
}

func thumbBirl() *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Pitch:  c.LowA,
			Length: c.Quarter,
			Embellishment: &emb.Embellishment{
				Type:    emb.Birl,
				Variant: emb.Thumb,
			},
		},
	}
}

func Test_birlsExpander_regular_ExpandSymbol(t *testing.T) {
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
			name: "regular",
			prepare: func(f *fields) {
				f.symbol = regBirl()
				f.want = []c.Pitch{c.LowA, c.LowG, c.LowA, c.LowG}
			},
		},
		{
			name: "regular with previous low g => 'half birl' ",
			prepare: func(f *fields) {
				f.symbol = regBirl()
				f.prevPitch = c.LowA
				f.want = []c.Pitch{c.LowG, c.LowA, c.LowG}
			},
		},
		{
			name: "g birl",
			prepare: func(f *fields) {
				f.symbol = gBirl()
				f.want = []c.Pitch{c.HighG, c.LowA, c.LowG, c.LowA, c.LowG}
			},
		},
		{
			name: "thumb birl",
			prepare: func(f *fields) {
				f.symbol = thumbBirl()
				f.want = []c.Pitch{c.HighA, c.LowA, c.LowG, c.LowA, c.LowG}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewBirlsExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}
