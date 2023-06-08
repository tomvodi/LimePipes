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

func regPele(pitch c.Pitch) *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Pitch:  pitch,
			Length: c.Quarter,
			Embellishment: &emb.Embellishment{
				Type: emb.Pele,
			},
		},
	}
}

func halfPele(pitch c.Pitch) *music_model.Symbol {
	return peleVar(pitch, emb.Half)
}

func thumbPele(pitch c.Pitch) *music_model.Symbol {
	return peleVar(pitch, emb.Thumb)
}

func peleVar(pitch c.Pitch, variant emb.EmbellishmentVariant) *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Pitch:  pitch,
			Length: c.Quarter,
			Embellishment: &emb.Embellishment{
				Type:    emb.Pele,
				Variant: variant,
			},
		},
	}
}

func makeLight(sym *music_model.Symbol) *music_model.Symbol {
	sym.Note.Embellishment.Weight = emb.Light
	return sym
}

func Test_peleExpander_regular_ExpandSymbol(t *testing.T) {
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
			name: "Low A",
			prepare: func(f *fields) {
				f.symbol = regPele(c.LowA)
				f.want = []c.Pitch{c.HighG, c.LowA, c.E, c.LowA, c.LowG}
			},
		},
		{
			name: "B",
			prepare: func(f *fields) {
				f.symbol = regPele(c.B)
				f.want = []c.Pitch{c.HighG, c.B, c.E, c.B, c.LowG}
			},
		},
		{
			name: "C",
			prepare: func(f *fields) {
				f.symbol = regPele(c.C)
				f.want = []c.Pitch{c.HighG, c.C, c.E, c.C, c.LowG}
			},
		},
		{
			name: "D",
			prepare: func(f *fields) {
				f.symbol = regPele(c.D)
				f.want = []c.Pitch{c.HighG, c.D, c.E, c.D, c.LowG}
			},
		},
		{
			name: "D light",
			prepare: func(f *fields) {
				f.symbol = makeLight(regPele(c.D))
				f.want = []c.Pitch{c.HighG, c.D, c.E, c.D, c.C}
			},
		},
		{
			name: "E",
			prepare: func(f *fields) {
				f.symbol = regPele(c.E)
				f.want = []c.Pitch{c.HighG, c.E, c.F, c.E, c.LowA}
			},
		},
		{
			name: "F",
			prepare: func(f *fields) {
				f.symbol = regPele(c.F)
				f.want = []c.Pitch{c.HighG, c.F, c.HighG, c.F, c.E}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewPelesExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}

func Test_peleExpander_half_ExpandSymbol(t *testing.T) {
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
			name: "Low A",
			prepare: func(f *fields) {
				f.symbol = halfPele(c.LowA)
				f.want = []c.Pitch{c.LowA, c.E, c.LowA, c.LowG}
			},
		},
		{
			name: "B",
			prepare: func(f *fields) {
				f.symbol = halfPele(c.B)
				f.want = []c.Pitch{c.B, c.E, c.B, c.LowG}
			},
		},
		{
			name: "C",
			prepare: func(f *fields) {
				f.symbol = halfPele(c.C)
				f.want = []c.Pitch{c.C, c.E, c.C, c.LowG}
			},
		},
		{
			name: "D",
			prepare: func(f *fields) {
				f.symbol = halfPele(c.D)
				f.want = []c.Pitch{c.D, c.E, c.D, c.LowG}
			},
		},
		{
			name: "D light",
			prepare: func(f *fields) {
				f.symbol = makeLight(halfPele(c.D))
				f.want = []c.Pitch{c.D, c.E, c.D, c.C}
			},
		},
		{
			name: "E",
			prepare: func(f *fields) {
				f.symbol = halfPele(c.E)
				f.want = []c.Pitch{c.E, c.F, c.E, c.LowA}
			},
		},
		{
			name: "F",
			prepare: func(f *fields) {
				f.symbol = halfPele(c.F)
				f.want = []c.Pitch{c.F, c.HighG, c.F, c.E}
			},
		},
		{
			name: "High G",
			prepare: func(f *fields) {
				f.symbol = halfPele(c.HighG)
				f.want = []c.Pitch{c.HighG, c.HighA, c.HighG, c.F}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewPelesExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}

func Test_peleExpander_thumb_ExpandSymbol(t *testing.T) {
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
			name: "Low A",
			prepare: func(f *fields) {
				f.symbol = thumbPele(c.LowA)
				f.want = []c.Pitch{c.HighA, c.LowA, c.E, c.LowA, c.LowG}
			},
		},
		{
			name: "B",
			prepare: func(f *fields) {
				f.symbol = thumbPele(c.B)
				f.want = []c.Pitch{c.HighA, c.B, c.E, c.B, c.LowG}
			},
		},
		{
			name: "C",
			prepare: func(f *fields) {
				f.symbol = thumbPele(c.C)
				f.want = []c.Pitch{c.HighA, c.C, c.E, c.C, c.LowG}
			},
		},
		{
			name: "D",
			prepare: func(f *fields) {
				f.symbol = thumbPele(c.D)
				f.want = []c.Pitch{c.HighA, c.D, c.E, c.D, c.LowG}
			},
		},
		{
			name: "D light",
			prepare: func(f *fields) {
				f.symbol = makeLight(thumbPele(c.D))
				f.want = []c.Pitch{c.HighA, c.D, c.E, c.D, c.C}
			},
		},
		{
			name: "E",
			prepare: func(f *fields) {
				f.symbol = thumbPele(c.E)
				f.want = []c.Pitch{c.HighA, c.E, c.F, c.E, c.LowA}
			},
		},
		{
			name: "F",
			prepare: func(f *fields) {
				f.symbol = thumbPele(c.F)
				f.want = []c.Pitch{c.HighA, c.F, c.HighG, c.F, c.E}
			},
		},
		{
			name: "High G",
			prepare: func(f *fields) {
				f.symbol = thumbPele(c.HighG)
				f.want = []c.Pitch{c.HighA, c.HighG, c.HighA, c.HighG, c.F}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewPelesExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}
