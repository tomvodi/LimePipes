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

func regDoubling(pitch c.Pitch) *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Pitch:  pitch,
			Length: c.Quarter,
			Embellishment: &emb.Embellishment{
				Type: emb.Doubling,
			},
		},
	}
}

func thumbDoubling(pitch c.Pitch) *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Pitch:  pitch,
			Length: c.Quarter,
			Embellishment: &emb.Embellishment{
				Type:    emb.Doubling,
				Variant: emb.Thumb,
			},
		},
	}
}

func halfDoubling(pitch c.Pitch) *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Pitch:  pitch,
			Length: c.Quarter,
			Embellishment: &emb.Embellishment{
				Type:    emb.Doubling,
				Variant: emb.Half,
			},
		},
	}
}

func Test_dblExpander_regular_ExpandSymbol(t *testing.T) {
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
				f.symbol = regDoubling(c.LowG)
				f.want = []c.Pitch{c.HighG, c.LowG, c.D}
			},
		},
		{
			name: "low a",
			prepare: func(f *fields) {
				f.symbol = regDoubling(c.LowA)
				f.want = []c.Pitch{c.HighG, c.LowA, c.D}
			},
		},
		{
			name: "b",
			prepare: func(f *fields) {
				f.symbol = regDoubling(c.B)
				f.want = []c.Pitch{c.HighG, c.B, c.D}
			},
		},
		{
			name: "c",
			prepare: func(f *fields) {
				f.symbol = regDoubling(c.C)
				f.want = []c.Pitch{c.HighG, c.C, c.D}
			},
		},
		{
			name: "d",
			prepare: func(f *fields) {
				f.symbol = regDoubling(c.D)
				f.want = []c.Pitch{c.HighG, c.D, c.E}
			},
		},
		{
			name: "e",
			prepare: func(f *fields) {
				f.symbol = regDoubling(c.E)
				f.want = []c.Pitch{c.HighG, c.E, c.F}
			},
		},
		{
			name: "f",
			prepare: func(f *fields) {
				f.symbol = regDoubling(c.F)
				f.want = []c.Pitch{c.HighG, c.F, c.HighG}
			},
		},
		{
			name: "hg",
			prepare: func(f *fields) {
				f.symbol = regDoubling(c.HighG)
				f.want = []c.Pitch{c.HighG, c.F}
			},
		},
		{
			name: "ha",
			prepare: func(f *fields) {
				f.symbol = regDoubling(c.HighA)
				f.want = []c.Pitch{c.HighA, c.HighG}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewDoublingsExpander()
			pack.ExpandSymbol(f.symbol, c.NoPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}

func Test_dblUnapcker_thumb_ExpandSymbol(t *testing.T) {
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
				f.symbol = thumbDoubling(c.LowG)
				f.want = []c.Pitch{c.HighA, c.LowG, c.D}
			},
		},
		{
			name: "low a",
			prepare: func(f *fields) {
				f.symbol = thumbDoubling(c.LowA)
				f.want = []c.Pitch{c.HighA, c.LowA, c.D}
			},
		},
		{
			name: "b",
			prepare: func(f *fields) {
				f.symbol = thumbDoubling(c.B)
				f.want = []c.Pitch{c.HighA, c.B, c.D}
			},
		},
		{
			name: "c",
			prepare: func(f *fields) {
				f.symbol = thumbDoubling(c.C)
				f.want = []c.Pitch{c.HighA, c.C, c.D}
			},
		},
		{
			name: "d",
			prepare: func(f *fields) {
				f.symbol = thumbDoubling(c.D)
				f.want = []c.Pitch{c.HighA, c.D, c.E}
			},
		},
		{
			name: "e",
			prepare: func(f *fields) {
				f.symbol = thumbDoubling(c.E)
				f.want = []c.Pitch{c.HighA, c.E, c.F}
			},
		},
		{
			name: "f",
			prepare: func(f *fields) {
				f.symbol = thumbDoubling(c.F)
				f.want = []c.Pitch{c.HighA, c.F, c.HighG}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewDoublingsExpander()
			pack.ExpandSymbol(f.symbol, c.NoPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}

func Test_dblUnapcker_half_ExpandSymbol(t *testing.T) {
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
				f.symbol = halfDoubling(c.LowG)
				f.want = []c.Pitch{c.LowG, c.D}
			},
		},
		{
			name: "low a",
			prepare: func(f *fields) {
				f.symbol = halfDoubling(c.LowA)
				f.want = []c.Pitch{c.LowA, c.D}
			},
		},
		{
			name: "b",
			prepare: func(f *fields) {
				f.symbol = halfDoubling(c.B)
				f.want = []c.Pitch{c.B, c.D}
			},
		},
		{
			name: "c",
			prepare: func(f *fields) {
				f.symbol = halfDoubling(c.C)
				f.want = []c.Pitch{c.C, c.D}
			},
		},
		{
			name: "d",
			prepare: func(f *fields) {
				f.symbol = halfDoubling(c.D)
				f.want = []c.Pitch{c.D, c.E}
			},
		},
		{
			name: "e",
			prepare: func(f *fields) {
				f.symbol = halfDoubling(c.E)
				f.want = []c.Pitch{c.E, c.F}
			},
		},
		{
			name: "f",
			prepare: func(f *fields) {
				f.symbol = halfDoubling(c.F)
				f.want = []c.Pitch{c.F, c.HighG}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewDoublingsExpander()
			pack.ExpandSymbol(f.symbol, c.NoPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}
