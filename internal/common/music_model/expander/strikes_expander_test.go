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

func singleStrike(pitch c.Pitch) *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Pitch:  pitch + 1,
			Length: c.Quarter,
			Embellishment: &emb.Embellishment{
				Type:  emb.Strike,
				Pitch: pitch,
			},
		},
	}
}

func halfStrike(pitch c.Pitch) *music_model.Symbol {
	return strikeVariant(pitch, emb.Half)
}

func thumbStrike(pitch c.Pitch) *music_model.Symbol {
	return strikeVariant(pitch, emb.Thumb)
}

func gStrike(pitch c.Pitch) *music_model.Symbol {
	return strikeVariant(pitch, emb.G)
}

func strikeVariant(pitch c.Pitch, variant emb.EmbellishmentVariant) *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Pitch:  pitch,
			Length: c.Quarter,
			Embellishment: &emb.Embellishment{
				Type:    emb.Strike,
				Variant: variant,
			},
		},
	}
}

func Test_strikesExpander_single_ExpandSymbol(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		symbol *music_model.Symbol
		want   []c.Pitch
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: "low g",
			prepare: func(f *fields) {
				f.symbol = singleStrike(c.LowG)
				f.want = []c.Pitch{c.LowG}
			},
		},
		{
			name: "low a",
			prepare: func(f *fields) {
				f.symbol = singleStrike(c.LowA)
				f.want = []c.Pitch{c.LowA}
			},
		},
		{
			name: "b",
			prepare: func(f *fields) {
				f.symbol = singleStrike(c.B)
				f.want = []c.Pitch{c.B}
			},
		},
		{
			name: "c",
			prepare: func(f *fields) {
				f.symbol = singleStrike(c.C)
				f.want = []c.Pitch{c.C}
			},
		},
		{
			name: "d",
			prepare: func(f *fields) {
				f.symbol = singleStrike(c.D)
				f.want = []c.Pitch{c.D}
			},
		},
		{
			name: "e",
			prepare: func(f *fields) {
				f.symbol = singleStrike(c.E)
				f.want = []c.Pitch{c.E}
			},
		},
		{
			name: "f",
			prepare: func(f *fields) {
				f.symbol = singleStrike(c.F)
				f.want = []c.Pitch{c.F}
			},
		},
		{
			name: "hg",
			prepare: func(f *fields) {
				f.symbol = singleStrike(c.HighG)
				f.want = []c.Pitch{c.HighG}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewStrikesExpander()
			pack.ExpandSymbol(f.symbol, c.NoPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}

func Test_strikesExpander_half_ExpandSymbol(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		symbol *music_model.Symbol
		want   []c.Pitch
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: "low a",
			prepare: func(f *fields) {
				f.symbol = halfStrike(c.LowA)
				f.want = []c.Pitch{c.LowA, c.LowG}
			},
		},
		{
			name: "b",
			prepare: func(f *fields) {
				f.symbol = halfStrike(c.B)
				f.want = []c.Pitch{c.B, c.LowG}
			},
		},
		{
			name: "c",
			prepare: func(f *fields) {
				f.symbol = halfStrike(c.C)
				f.want = []c.Pitch{c.C, c.LowG}
			},
		},
		{
			name: "d",
			prepare: func(f *fields) {
				f.symbol = halfStrike(c.D)
				f.want = []c.Pitch{c.D, c.LowG}
			},
		},
		{
			name: "d light",
			prepare: func(f *fields) {
				str := halfStrike(c.D)
				str.Note.Embellishment.Weight = emb.Light
				f.symbol = str
				f.want = []c.Pitch{c.D, c.C}
			},
		},
		{
			name: "e",
			prepare: func(f *fields) {
				f.symbol = halfStrike(c.E)
				f.want = []c.Pitch{c.E, c.LowA}
			},
		},
		{
			name: "f",
			prepare: func(f *fields) {
				f.symbol = halfStrike(c.F)
				f.want = []c.Pitch{c.F, c.E}
			},
		},
		{
			name: "hg",
			prepare: func(f *fields) {
				f.symbol = halfStrike(c.HighG)
				f.want = []c.Pitch{c.HighG, c.F}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewStrikesExpander()
			pack.ExpandSymbol(f.symbol, c.NoPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}

func Test_strikesExpander_thumb_ExpandSymbol(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		symbol *music_model.Symbol
		want   []c.Pitch
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: "low a",
			prepare: func(f *fields) {
				f.symbol = thumbStrike(c.LowA)
				f.want = []c.Pitch{c.HighA, c.LowA, c.LowG}
			},
		},
		{
			name: "b",
			prepare: func(f *fields) {
				f.symbol = thumbStrike(c.B)
				f.want = []c.Pitch{c.HighA, c.B, c.LowG}
			},
		},
		{
			name: "c",
			prepare: func(f *fields) {
				f.symbol = thumbStrike(c.C)
				f.want = []c.Pitch{c.HighA, c.C, c.LowG}
			},
		},
		{
			name: "d",
			prepare: func(f *fields) {
				f.symbol = thumbStrike(c.D)
				f.want = []c.Pitch{c.HighA, c.D, c.LowG}
			},
		},
		{
			name: "d light",
			prepare: func(f *fields) {
				str := thumbStrike(c.D)
				str.Note.Embellishment.Weight = emb.Light
				f.symbol = str
				f.want = []c.Pitch{c.HighA, c.D, c.C}
			},
		},
		{
			name: "e",
			prepare: func(f *fields) {
				f.symbol = thumbStrike(c.E)
				f.want = []c.Pitch{c.HighA, c.E, c.LowA}
			},
		},
		{
			name: "f",
			prepare: func(f *fields) {
				f.symbol = thumbStrike(c.F)
				f.want = []c.Pitch{c.HighA, c.F, c.E}
			},
		},
		{
			name: "hg",
			prepare: func(f *fields) {
				f.symbol = thumbStrike(c.HighG)
				f.want = []c.Pitch{c.HighA, c.HighG, c.F}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewStrikesExpander()
			pack.ExpandSymbol(f.symbol, c.NoPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}

func Test_strikesExpander_g_ExpandSymbol(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		symbol *music_model.Symbol
		want   []c.Pitch
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: "low a",
			prepare: func(f *fields) {
				f.symbol = gStrike(c.LowA)
				f.want = []c.Pitch{c.HighG, c.LowA, c.LowG}
			},
		},
		{
			name: "b",
			prepare: func(f *fields) {
				f.symbol = gStrike(c.B)
				f.want = []c.Pitch{c.HighG, c.B, c.LowG}
			},
		},
		{
			name: "c",
			prepare: func(f *fields) {
				f.symbol = gStrike(c.C)
				f.want = []c.Pitch{c.HighG, c.C, c.LowG}
			},
		},
		{
			name: "d",
			prepare: func(f *fields) {
				f.symbol = gStrike(c.D)
				f.want = []c.Pitch{c.HighG, c.D, c.LowG}
			},
		},
		{
			name: "d light",
			prepare: func(f *fields) {
				str := gStrike(c.D)
				str.Note.Embellishment.Weight = emb.Light
				f.symbol = str
				f.want = []c.Pitch{c.HighG, c.D, c.C}
			},
		},
		{
			name: "e",
			prepare: func(f *fields) {
				f.symbol = gStrike(c.E)
				f.want = []c.Pitch{c.HighG, c.E, c.LowA}
			},
		},
		{
			name: "f",
			prepare: func(f *fields) {
				f.symbol = gStrike(c.F)
				f.want = []c.Pitch{c.HighG, c.F, c.E}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewStrikesExpander()
			pack.ExpandSymbol(f.symbol, c.NoPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}
