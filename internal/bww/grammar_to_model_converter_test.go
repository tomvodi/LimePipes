package bww

import (
	. "github.com/onsi/gomega"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/boundary"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/length"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/measure"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/pitch"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/tuplet"
	"github.com/tomvodi/limepipes/internal/utils"
	"testing"
)

func Test_handleTriplet(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		measure *measure.Measure
		sym     string
		wantErr bool
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		after   func(f *fields)
	}{
		{
			name: "no symbols in measure",
			prepare: func(f *fields) {
				f.measure = &measure.Measure{}
				f.wantErr = true
			},
		},
		{
			name: "not enough symbols in measure",
			prepare: func(f *fields) {
				f.measure = &measure.Measure{
					Time: nil,
					Symbols: []*symbols.Symbol{
						{Note: &symbols.Note{Pitch: pitch.Pitch_LowA, Length: length.Length_Eighth}},
						{Note: &symbols.Note{Pitch: pitch.Pitch_LowA, Length: length.Length_Eighth}},
					},
				}
				f.wantErr = true
			},
		},
		{
			name: "not all preceding symbols are notes",
			prepare: func(f *fields) {
				f.measure = &measure.Measure{
					Time: nil,
					Symbols: []*symbols.Symbol{
						{Note: &symbols.Note{Pitch: pitch.Pitch_LowA, Length: length.Length_Eighth}},
						{Note: &symbols.Note{
							Embellishment: &embellishment.Embellishment{Type: embellishment.Type_Doubling},
						}},
						{Note: &symbols.Note{Pitch: pitch.Pitch_LowA, Length: length.Length_Eighth}},
					},
				}
				f.wantErr = true
			},
		},
		{
			name: "all preceding symbols are notes",
			prepare: func(f *fields) {
				f.measure = &measure.Measure{
					Time: nil,
					Symbols: []*symbols.Symbol{
						{Note: &symbols.Note{Pitch: pitch.Pitch_LowA, Length: length.Length_Eighth}},
						{Note: &symbols.Note{Pitch: pitch.Pitch_LowA, Length: length.Length_Eighth}},
						{Note: &symbols.Note{Pitch: pitch.Pitch_LowA, Length: length.Length_Eighth}},
					},
				}
				f.wantErr = false
			},
			after: func(f *fields) {
				g.Expect(f.measure.Symbols).To(HaveLen(5))
				g.Expect(f.measure.Symbols[0].Tuplet).To(Equal(&tuplet.Tuplet{
					BoundaryType: boundary.Boundary_Start,
					VisibleNotes: 3,
					PlayedNotes:  2,
				}))
				g.Expect(f.measure.Symbols[4].Tuplet).To(Equal(&tuplet.Tuplet{
					BoundaryType: boundary.Boundary_End,
					VisibleNotes: 3,
					PlayedNotes:  2,
				}))
			},
		},
		{
			name: "if there is already a tuplet start, don't add another one",
			prepare: func(f *fields) {
				f.measure = &measure.Measure{
					Time: nil,
					Symbols: []*symbols.Symbol{
						{
							Tuplet: &tuplet.Tuplet{
								BoundaryType: boundary.Boundary_Start,
								VisibleNotes: 3,
								PlayedNotes:  2,
							},
						},
						{Note: &symbols.Note{Pitch: pitch.Pitch_LowA, Length: length.Length_Eighth}},
						{Note: &symbols.Note{Pitch: pitch.Pitch_LowA, Length: length.Length_Eighth}},
						{Note: &symbols.Note{Pitch: pitch.Pitch_LowA, Length: length.Length_Eighth}},
					},
				}
				f.wantErr = false
			},
			after: func(f *fields) {
				g.Expect(f.measure.Symbols).To(HaveLen(5))
				g.Expect(f.measure.Symbols[0].Tuplet).To(Equal(&tuplet.Tuplet{
					BoundaryType: boundary.Boundary_Start,
					VisibleNotes: 3,
					PlayedNotes:  2,
				}))
				g.Expect(f.measure.Symbols[4].Tuplet).To(Equal(&tuplet.Tuplet{
					BoundaryType: boundary.Boundary_End,
					VisibleNotes: 3,
					PlayedNotes:  2,
				}))
			},
		},
		{
			name: "if there is a tuplet end, a start mus be added",
			prepare: func(f *fields) {
				f.measure = &measure.Measure{
					Time: nil,
					Symbols: []*symbols.Symbol{
						{
							Tuplet: &tuplet.Tuplet{
								BoundaryType: boundary.Boundary_End,
								VisibleNotes: 7,
								PlayedNotes:  6,
							},
						},
						{Note: &symbols.Note{Pitch: pitch.Pitch_LowA, Length: length.Length_Eighth}},
						{Note: &symbols.Note{Pitch: pitch.Pitch_LowA, Length: length.Length_Eighth}},
						{Note: &symbols.Note{Pitch: pitch.Pitch_LowA, Length: length.Length_Eighth}},
					},
				}
				f.wantErr = false
			},
			after: func(f *fields) {
				g.Expect(f.measure.Symbols).To(HaveLen(6))
				g.Expect(f.measure.Symbols[0].Tuplet).To(Equal(&tuplet.Tuplet{
					BoundaryType: boundary.Boundary_End,
					VisibleNotes: 7,
					PlayedNotes:  6,
				}))
				g.Expect(f.measure.Symbols[1].Tuplet).To(Equal(&tuplet.Tuplet{
					BoundaryType: boundary.Boundary_Start,
					VisibleNotes: 3,
					PlayedNotes:  2,
				}))
				g.Expect(f.measure.Symbols[5].Tuplet).To(Equal(&tuplet.Tuplet{
					BoundaryType: boundary.Boundary_End,
					VisibleNotes: 3,
					PlayedNotes:  2,
				}))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			err := handleTriplet(f.measure, f.sym)
			if f.wantErr {
				g.Expect(err).Should(HaveOccurred())
			} else {
				g.Expect(err).ShouldNot(HaveOccurred())
			}

			if tt.after != nil {
				tt.after(f)
			}
		})
	}
}
func Test_pitchesFromFermataSym(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		sym     string
		fermata bool
		want    []pitch.Pitch
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: "cadence ged",
			prepare: func(f *fields) {
				f.sym = "cadged"
				f.fermata = false
				f.want = []pitch.Pitch{
					pitch.Pitch_HighG,
					pitch.Pitch_E,
					pitch.Pitch_D,
				}
			},
		},
		{
			name: "fermata cadence ged",
			prepare: func(f *fields) {
				f.sym = "fcadged"
				f.fermata = true
				f.want = []pitch.Pitch{
					pitch.Pitch_HighG,
					pitch.Pitch_E,
					pitch.Pitch_D,
				}
			},
		},
		{
			name: "fermata cadence af",
			prepare: func(f *fields) {
				f.sym = "fcadaf"
				f.fermata = true
				f.want = []pitch.Pitch{
					pitch.Pitch_HighA,
					pitch.Pitch_F,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			got := pitchesFromCadenceSym(f.sym, f.fermata)
			g.Expect(got).To(Equal(f.want))

		})
	}
}
