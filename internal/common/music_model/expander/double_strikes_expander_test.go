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

func regDoubleStrike(pitch c.Pitch) *music_model.Symbol {
	return doubleStrikeVar(pitch, emb.NoVariant)
}

func halfDoubleStrike(pitch c.Pitch) *music_model.Symbol {
	return doubleStrikeVar(pitch, emb.Half)
}

func thumbDoubleStrike(pitch c.Pitch) *music_model.Symbol {
	return doubleStrikeVar(pitch, emb.Thumb)
}

func gDoubleStrike(pitch c.Pitch) *music_model.Symbol {
	return doubleStrikeVar(pitch, emb.G)
}

func doubleStrikeVar(pitch c.Pitch, variant emb.EmbellishmentVariant) *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Pitch:  pitch,
			Length: c.Quarter,
			Embellishment: &emb.Embellishment{
				Type:    emb.DoubleStrike,
				Variant: variant,
			},
		},
	}
}

func Test_doubleStrikeExpander_regular_ExpandSymbol(t *testing.T) {
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
				f.symbol = regDoubleStrike(c.LowA)
				f.want = []c.Pitch{c.LowG, c.LowA, c.LowG}
			},
		},
		{
			name: "B",
			prepare: func(f *fields) {
				f.symbol = regDoubleStrike(c.B)
				f.want = []c.Pitch{c.LowG, c.B, c.LowG}
			},
		},
		{
			name: "C",
			prepare: func(f *fields) {
				f.symbol = regDoubleStrike(c.C)
				f.want = []c.Pitch{c.LowG, c.C, c.LowG}
			},
		},
		{
			name: "D",
			prepare: func(f *fields) {
				f.symbol = regDoubleStrike(c.D)
				f.want = []c.Pitch{c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "D light",
			prepare: func(f *fields) {
				f.symbol = makeLight(regDoubleStrike(c.D))
				f.want = []c.Pitch{c.C, c.D, c.C}
			},
		},
		{
			name: "E",
			prepare: func(f *fields) {
				f.symbol = regDoubleStrike(c.E)
				f.want = []c.Pitch{c.LowA, c.E, c.LowA}
			},
		},
		{
			name: "F",
			prepare: func(f *fields) {
				f.symbol = regDoubleStrike(c.F)
				f.want = []c.Pitch{c.E, c.F, c.E}
			},
		},
		{
			name: "High G",
			prepare: func(f *fields) {
				f.symbol = regDoubleStrike(c.HighG)
				f.want = []c.Pitch{c.F, c.HighG, c.F}
			},
		},
		{
			name: "High A",
			prepare: func(f *fields) {
				f.symbol = regDoubleStrike(c.HighA)
				f.want = []c.Pitch{c.HighG, c.HighA, c.HighG}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewDoubleStrikesExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}

func Test_doubleStrikeExpander_g_ExpandSymbol(t *testing.T) {
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
				f.symbol = gDoubleStrike(c.LowA)
				f.want = []c.Pitch{c.HighG, c.LowA, c.LowG, c.LowA, c.LowG}
			},
		},
		{
			name: "B",
			prepare: func(f *fields) {
				f.symbol = gDoubleStrike(c.B)
				f.want = []c.Pitch{c.HighG, c.B, c.LowG, c.B, c.LowG}
			},
		},
		{
			name: "C",
			prepare: func(f *fields) {
				f.symbol = gDoubleStrike(c.C)
				f.want = []c.Pitch{c.HighG, c.C, c.LowG, c.C, c.LowG}
			},
		},
		{
			name: "D",
			prepare: func(f *fields) {
				f.symbol = gDoubleStrike(c.D)
				f.want = []c.Pitch{c.HighG, c.D, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "D light",
			prepare: func(f *fields) {
				f.symbol = makeLight(gDoubleStrike(c.D))
				f.want = []c.Pitch{c.HighG, c.D, c.C, c.D, c.C}
			},
		},
		{
			name: "E",
			prepare: func(f *fields) {
				f.symbol = gDoubleStrike(c.E)
				f.want = []c.Pitch{c.HighG, c.E, c.LowA, c.E, c.LowA}
			},
		},
		{
			name: "F",
			prepare: func(f *fields) {
				f.symbol = gDoubleStrike(c.F)
				f.want = []c.Pitch{c.HighG, c.F, c.E, c.F, c.E}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewDoubleStrikesExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}

func Test_doubleStrikeExpander_thumb_ExpandSymbol(t *testing.T) {
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
				f.symbol = thumbDoubleStrike(c.LowA)
				f.want = []c.Pitch{c.HighA, c.LowA, c.LowG, c.LowA, c.LowG}
			},
		},
		{
			name: "B",
			prepare: func(f *fields) {
				f.symbol = thumbDoubleStrike(c.B)
				f.want = []c.Pitch{c.HighA, c.B, c.LowG, c.B, c.LowG}
			},
		},
		{
			name: "C",
			prepare: func(f *fields) {
				f.symbol = thumbDoubleStrike(c.C)
				f.want = []c.Pitch{c.HighA, c.C, c.LowG, c.C, c.LowG}
			},
		},
		{
			name: "D",
			prepare: func(f *fields) {
				f.symbol = thumbDoubleStrike(c.D)
				f.want = []c.Pitch{c.HighA, c.D, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "D light",
			prepare: func(f *fields) {
				f.symbol = makeLight(thumbDoubleStrike(c.D))
				f.want = []c.Pitch{c.HighA, c.D, c.C, c.D, c.C}
			},
		},
		{
			name: "E",
			prepare: func(f *fields) {
				f.symbol = thumbDoubleStrike(c.E)
				f.want = []c.Pitch{c.HighA, c.E, c.LowA, c.E, c.LowA}
			},
		},
		{
			name: "F",
			prepare: func(f *fields) {
				f.symbol = thumbDoubleStrike(c.F)
				f.want = []c.Pitch{c.HighA, c.F, c.E, c.F, c.E}
			},
		},
		{
			name: "High G",
			prepare: func(f *fields) {
				f.symbol = thumbDoubleStrike(c.HighG)
				f.want = []c.Pitch{c.HighA, c.HighG, c.F, c.HighG, c.F}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewDoubleStrikesExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}

func Test_doubleStrikeExpander_half_ExpandSymbol(t *testing.T) {
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
				f.symbol = halfDoubleStrike(c.LowA)
				f.want = []c.Pitch{c.LowA, c.LowG, c.LowA, c.LowG}
			},
		},
		{
			name: "B",
			prepare: func(f *fields) {
				f.symbol = halfDoubleStrike(c.B)
				f.want = []c.Pitch{c.B, c.LowG, c.B, c.LowG}
			},
		},
		{
			name: "C",
			prepare: func(f *fields) {
				f.symbol = halfDoubleStrike(c.C)
				f.want = []c.Pitch{c.C, c.LowG, c.C, c.LowG}
			},
		},
		{
			name: "D",
			prepare: func(f *fields) {
				f.symbol = halfDoubleStrike(c.D)
				f.want = []c.Pitch{c.D, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "D light",
			prepare: func(f *fields) {
				f.symbol = makeLight(halfDoubleStrike(c.D))
				f.want = []c.Pitch{c.D, c.C, c.D, c.C}
			},
		},
		{
			name: "E",
			prepare: func(f *fields) {
				f.symbol = halfDoubleStrike(c.E)
				f.want = []c.Pitch{c.E, c.LowA, c.E, c.LowA}
			},
		},
		{
			name: "F",
			prepare: func(f *fields) {
				f.symbol = halfDoubleStrike(c.F)
				f.want = []c.Pitch{c.F, c.E, c.F, c.E}
			},
		},
		{
			name: "High G",
			prepare: func(f *fields) {
				f.symbol = halfDoubleStrike(c.HighG)
				f.want = []c.Pitch{c.HighG, c.F, c.HighG, c.F}
			},
		},
		{
			name: "High A",
			prepare: func(f *fields) {
				f.symbol = halfDoubleStrike(c.HighA)
				f.want = []c.Pitch{c.HighA, c.HighG, c.HighA, c.HighG}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewDoubleStrikesExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}
