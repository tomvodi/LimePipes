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

func regGrp() *music_model.Symbol {
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

func gGrp(pitch c.Pitch) *music_model.Symbol {
	return gripVariant(pitch, emb.G)
}

func thumbGrip(pitch c.Pitch) *music_model.Symbol {
	return gripVariant(pitch, emb.Thumb)
}

func halfGrip(pitch c.Pitch) *music_model.Symbol {
	return gripVariant(pitch, emb.Half)
}

func gripVariant(pitch c.Pitch, variant emb.EmbellishmentVariant) *music_model.Symbol {
	return &music_model.Symbol{
		Note: &symbols.Note{
			Pitch:  pitch,
			Length: c.Quarter,
			Embellishment: &emb.Embellishment{
				Type:    emb.Grip,
				Variant: variant,
			},
		},
	}
}

func makeB(sym *music_model.Symbol) *music_model.Symbol {
	sym.Note.Embellishment.Pitch = c.B
	return sym
}

func Test_grpExpander_regular_ExpandSymbol(t *testing.T) {
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
				f.symbol = regGrp()
				f.want = []c.Pitch{c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "regular with previous low g => half grip",
			prepare: func(f *fields) {
				f.symbol = regGrp()
				f.prevPitch = c.LowG
				f.want = []c.Pitch{c.D, c.LowG}
			},
		},
		{
			name: "regular b",
			prepare: func(f *fields) {
				f.symbol = makeB(regGrp())
				f.want = []c.Pitch{c.LowG, c.B, c.LowG}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewGripsExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}

func Test_grpExpander_g_ExpandSymbol(t *testing.T) {
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
			name: "LowA",
			prepare: func(f *fields) {
				f.symbol = gGrp(c.LowA)
				f.want = []c.Pitch{c.HighG, c.LowA, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "B",
			prepare: func(f *fields) {
				f.symbol = gGrp(c.B)
				f.want = []c.Pitch{c.HighG, c.B, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "c",
			prepare: func(f *fields) {
				f.symbol = gGrp(c.C)
				f.want = []c.Pitch{c.HighG, c.C, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "D",
			prepare: func(f *fields) {
				f.symbol = gGrp(c.D)
				f.want = []c.Pitch{c.HighG, c.D, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "D - B",
			prepare: func(f *fields) {
				f.symbol = makeB(gGrp(c.D))
				f.want = []c.Pitch{c.HighG, c.D, c.LowG, c.B, c.LowG}
			},
		},
		{
			name: "E",
			prepare: func(f *fields) {
				f.symbol = gGrp(c.E)
				f.want = []c.Pitch{c.HighG, c.E, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "F",
			prepare: func(f *fields) {
				f.symbol = gGrp(c.F)
				f.want = []c.Pitch{c.HighG, c.F, c.LowG, c.F, c.LowG}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewGripsExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}

func Test_grpExpander_thumb_ExpandSymbol(t *testing.T) {
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
			name: "LowA",
			prepare: func(f *fields) {
				f.symbol = thumbGrip(c.LowA)
				f.want = []c.Pitch{c.HighA, c.LowA, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "B",
			prepare: func(f *fields) {
				f.symbol = thumbGrip(c.B)
				f.want = []c.Pitch{c.HighA, c.B, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "c",
			prepare: func(f *fields) {
				f.symbol = thumbGrip(c.C)
				f.want = []c.Pitch{c.HighA, c.C, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "D",
			prepare: func(f *fields) {
				f.symbol = thumbGrip(c.D)
				f.want = []c.Pitch{c.HighA, c.D, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "D - B",
			prepare: func(f *fields) {
				f.symbol = makeB(thumbGrip(c.D))
				f.want = []c.Pitch{c.HighA, c.D, c.LowG, c.B, c.LowG}
			},
		},
		{
			name: "E",
			prepare: func(f *fields) {
				f.symbol = thumbGrip(c.E)
				f.want = []c.Pitch{c.HighA, c.E, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "F",
			prepare: func(f *fields) {
				f.symbol = thumbGrip(c.F)
				f.want = []c.Pitch{c.HighA, c.F, c.LowG, c.F, c.LowG}
			},
		},
		{
			name: "High G",
			prepare: func(f *fields) {
				f.symbol = thumbGrip(c.HighG)
				f.want = []c.Pitch{c.HighA, c.HighG, c.LowG, c.F, c.LowG}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewGripsExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}

func Test_grpExpander_half_ExpandSymbol(t *testing.T) {
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
			name: "LowA",
			prepare: func(f *fields) {
				f.symbol = halfGrip(c.LowA)
				f.want = []c.Pitch{c.LowA, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "B",
			prepare: func(f *fields) {
				f.symbol = halfGrip(c.B)
				f.want = []c.Pitch{c.B, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "c",
			prepare: func(f *fields) {
				f.symbol = halfGrip(c.C)
				f.want = []c.Pitch{c.C, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "D",
			prepare: func(f *fields) {
				f.symbol = halfGrip(c.D)
				f.want = []c.Pitch{c.D, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "D - B",
			prepare: func(f *fields) {
				f.symbol = makeB(halfGrip(c.D))
				f.want = []c.Pitch{c.D, c.LowG, c.B, c.LowG}
			},
		},
		{
			name: "E",
			prepare: func(f *fields) {
				f.symbol = halfGrip(c.E)
				f.want = []c.Pitch{c.E, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "F",
			prepare: func(f *fields) {
				f.symbol = halfGrip(c.F)
				f.want = []c.Pitch{c.F, c.LowG, c.F, c.LowG}
			},
		},
		{
			name: "High G",
			prepare: func(f *fields) {
				f.symbol = halfGrip(c.HighG)
				f.want = []c.Pitch{c.HighG, c.LowG, c.D, c.LowG}
			},
		},
		{
			name: "High A",
			prepare: func(f *fields) {
				f.symbol = halfGrip(c.HighA)
				f.want = []c.Pitch{c.HighA, c.LowG, c.D, c.LowG}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			pack := NewGripsExpander()
			pack.ExpandSymbol(f.symbol, f.prevPitch)
			want := fmt.Sprintf("%v", f.want)
			got := fmt.Sprintf("%v", f.symbol.Note.ExpandedEmbellishment)
			g.Expect(got).To(Equal(want))
		})
	}
}
