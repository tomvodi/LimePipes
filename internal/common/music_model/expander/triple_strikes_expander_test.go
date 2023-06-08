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

func regTripleStrike(pitch c.Pitch) *music_model.Symbol {
	return tripleStrikeVar(pitch, emb.NoVariant)
}

func halfTripleStrike(pitch c.Pitch) *music_model.Symbol {
	return tripleStrikeVar(pitch, emb.Half)
}

func thumbTripleStrike(pitch c.Pitch) *music_model.Symbol {
	return tripleStrikeVar(pitch, emb.Thumb)
}

func gTripleStrike(pitch c.Pitch) *music_model.Symbol {
	return tripleStrikeVar(pitch, emb.G)
}

func tripleStrikeVar(pitch c.Pitch, variant emb.EmbellishmentVariant) *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Pitch:  pitch,
			Length: c.Quarter,
			Embellishment: &emb.Embellishment{
				Type:    emb.TripleStrike,
				Variant: variant,
			},
		},
	}
}

func Test_tripleStrikeExpander_regular_ExpandSymbol(t *testing.T) {
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
				f.symbol = regTripleStrike(c.LowA)
				f.want = []c.Pitch{c.LowG, c.LowA, c.LowG, c.LowA, c.LowG}
			},
		},
		{
			name: "B",
			prepare: func(f *fields) {
				f.symbol = regTripleStrike(c.B)
				f.want = []c.Pitch{c.LowG, c.B, c.LowG, c.B, c.LowG}
			},
		},
		{
			name: "C",
			prepare: func(f *fields) {
				f.symbol = regTripleStrike(c.C)
				f.want = []c.Pitch{c.LowG, c.C, c.LowG, c.C, c.LowG}
			},
		},
		{
			name: "D",
			prepare: func(f *fields) {
				f.symbol = regTripleStrike(c.D)
				f.want = []c.Pitch{c.LowG, c.D, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "D light",
			prepare: func(f *fields) {
				f.symbol = makeLight(regTripleStrike(c.D))
				f.want = []c.Pitch{c.C, c.D, c.C, c.D, c.C}
			},
		},
		{
			name: "E",
			prepare: func(f *fields) {
				f.symbol = regTripleStrike(c.E)
				f.want = []c.Pitch{c.LowA, c.E, c.LowA, c.E, c.LowA}
			},
		},
		{
			name: "F",
			prepare: func(f *fields) {
				f.symbol = regTripleStrike(c.F)
				f.want = []c.Pitch{c.E, c.F, c.E, c.F, c.E}
			},
		},
		{
			name: "High G",
			prepare: func(f *fields) {
				f.symbol = regTripleStrike(c.HighG)
				f.want = []c.Pitch{c.F, c.HighG, c.F, c.HighG, c.F}
			},
		},
		{
			name: "High A",
			prepare: func(f *fields) {
				f.symbol = regTripleStrike(c.HighA)
				f.want = []c.Pitch{c.HighG, c.HighA, c.HighG, c.HighA, c.HighG}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewTripleStrikesExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}

func Test_tripleStrikeExpander_g_ExpandSymbol(t *testing.T) {
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
				f.symbol = gTripleStrike(c.LowA)
				f.want = []c.Pitch{c.HighG, c.LowA, c.LowG, c.LowA, c.LowG, c.LowA, c.LowG}
			},
		},
		{
			name: "B",
			prepare: func(f *fields) {
				f.symbol = gTripleStrike(c.B)
				f.want = []c.Pitch{c.HighG, c.B, c.LowG, c.B, c.LowG, c.B, c.LowG}
			},
		},
		{
			name: "C",
			prepare: func(f *fields) {
				f.symbol = gTripleStrike(c.C)
				f.want = []c.Pitch{c.HighG, c.C, c.LowG, c.C, c.LowG, c.C, c.LowG}
			},
		},
		{
			name: "D",
			prepare: func(f *fields) {
				f.symbol = gTripleStrike(c.D)
				f.want = []c.Pitch{c.HighG, c.D, c.LowG, c.D, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "D light",
			prepare: func(f *fields) {
				f.symbol = makeLight(gTripleStrike(c.D))
				f.want = []c.Pitch{c.HighG, c.D, c.C, c.D, c.C, c.D, c.C}
			},
		},
		{
			name: "E",
			prepare: func(f *fields) {
				f.symbol = gTripleStrike(c.E)
				f.want = []c.Pitch{c.HighG, c.E, c.LowA, c.E, c.LowA, c.E, c.LowA}
			},
		},
		{
			name: "F",
			prepare: func(f *fields) {
				f.symbol = gTripleStrike(c.F)
				f.want = []c.Pitch{c.HighG, c.F, c.E, c.F, c.E, c.F, c.E}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewTripleStrikesExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}

func Test_tripleStrikeExpander_thumb_ExpandSymbol(t *testing.T) {
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
				f.symbol = thumbTripleStrike(c.LowA)
				f.want = []c.Pitch{c.HighA, c.LowA, c.LowG, c.LowA, c.LowG, c.LowA, c.LowG}
			},
		},
		{
			name: "B",
			prepare: func(f *fields) {
				f.symbol = thumbTripleStrike(c.B)
				f.want = []c.Pitch{c.HighA, c.B, c.LowG, c.B, c.LowG, c.B, c.LowG}
			},
		},
		{
			name: "C",
			prepare: func(f *fields) {
				f.symbol = thumbTripleStrike(c.C)
				f.want = []c.Pitch{c.HighA, c.C, c.LowG, c.C, c.LowG, c.C, c.LowG}
			},
		},
		{
			name: "D",
			prepare: func(f *fields) {
				f.symbol = thumbTripleStrike(c.D)
				f.want = []c.Pitch{c.HighA, c.D, c.LowG, c.D, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "D light",
			prepare: func(f *fields) {
				f.symbol = makeLight(thumbTripleStrike(c.D))
				f.want = []c.Pitch{c.HighA, c.D, c.C, c.D, c.C, c.D, c.C}
			},
		},
		{
			name: "E",
			prepare: func(f *fields) {
				f.symbol = thumbTripleStrike(c.E)
				f.want = []c.Pitch{c.HighA, c.E, c.LowA, c.E, c.LowA, c.E, c.LowA}
			},
		},
		{
			name: "F",
			prepare: func(f *fields) {
				f.symbol = thumbTripleStrike(c.F)
				f.want = []c.Pitch{c.HighA, c.F, c.E, c.F, c.E, c.F, c.E}
			},
		},
		{
			name: "High G",
			prepare: func(f *fields) {
				f.symbol = thumbTripleStrike(c.HighG)
				f.want = []c.Pitch{c.HighA, c.HighG, c.F, c.HighG, c.F, c.HighG, c.F}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewTripleStrikesExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}

func Test_tripleStrikeExpander_half_ExpandSymbol(t *testing.T) {
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
				f.symbol = halfTripleStrike(c.LowA)
				f.want = []c.Pitch{c.LowA, c.LowG, c.LowA, c.LowG, c.LowA, c.LowG}
			},
		},
		{
			name: "B",
			prepare: func(f *fields) {
				f.symbol = halfTripleStrike(c.B)
				f.want = []c.Pitch{c.B, c.LowG, c.B, c.LowG, c.B, c.LowG}
			},
		},
		{
			name: "C",
			prepare: func(f *fields) {
				f.symbol = halfTripleStrike(c.C)
				f.want = []c.Pitch{c.C, c.LowG, c.C, c.LowG, c.C, c.LowG}
			},
		},
		{
			name: "D",
			prepare: func(f *fields) {
				f.symbol = halfTripleStrike(c.D)
				f.want = []c.Pitch{c.D, c.LowG, c.D, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "D light",
			prepare: func(f *fields) {
				f.symbol = makeLight(halfTripleStrike(c.D))
				f.want = []c.Pitch{c.D, c.C, c.D, c.C, c.D, c.C}
			},
		},
		{
			name: "E",
			prepare: func(f *fields) {
				f.symbol = halfTripleStrike(c.E)
				f.want = []c.Pitch{c.E, c.LowA, c.E, c.LowA, c.E, c.LowA}
			},
		},
		{
			name: "F",
			prepare: func(f *fields) {
				f.symbol = halfTripleStrike(c.F)
				f.want = []c.Pitch{c.F, c.E, c.F, c.E, c.F, c.E}
			},
		},
		{
			name: "High G",
			prepare: func(f *fields) {
				f.symbol = halfTripleStrike(c.HighG)
				f.want = []c.Pitch{c.HighG, c.F, c.HighG, c.F, c.HighG, c.F}
			},
		},
		{
			name: "High A",
			prepare: func(f *fields) {
				f.symbol = halfTripleStrike(c.HighA)
				f.want = []c.Pitch{c.HighA, c.HighG, c.HighA, c.HighG, c.HighA, c.HighG}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewTripleStrikesExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}
