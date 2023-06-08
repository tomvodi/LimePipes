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

func regBubbly() *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Pitch:  c.C,
			Length: c.Quarter,
			Embellishment: &emb.Embellishment{
				Type: emb.Grip,
			},
		},
	}
}

func Test_bubblyExpander_regular_ExpandSymbol(t *testing.T) {
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
				f.symbol = regBubbly()
				f.want = []c.Pitch{c.LowG, c.D, c.LowG, c.C, c.LowG}
			},
		},
		{
			name: "regular with previous low g => half bubbly",
			prepare: func(f *fields) {
				f.symbol = regBubbly()
				f.prevPitch = c.LowG
				f.want = []c.Pitch{c.D, c.LowG, c.C, c.LowG}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewBubblysExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}
