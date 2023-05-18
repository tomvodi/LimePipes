package bww

import (
	"banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/symbols"
	"banduslib/internal/common/music_model/symbols/tuplet"
	"banduslib/internal/utils"
	. "github.com/onsi/gomega"
	"testing"
)

func Test_handleTriplet(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		measure *music_model.Measure
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
				f.measure = &music_model.Measure{}
				f.wantErr = true
			},
		},
		{
			name: "not enough symbols in measure",
			prepare: func(f *fields) {
				f.measure = &music_model.Measure{
					Time: nil,
					Symbols: []*music_model.Symbol{
						{Note: &symbols.Note{Pitch: common.LowA, Length: common.Eighth}},
						{Note: &symbols.Note{Pitch: common.LowA, Length: common.Eighth}},
					},
				}
				f.wantErr = true
			},
		},
		{
			name: "not all preceding symbols are notes",
			prepare: func(f *fields) {
				f.measure = &music_model.Measure{
					Time: nil,
					Symbols: []*music_model.Symbol{
						{Note: &symbols.Note{Pitch: common.LowA, Length: common.Eighth}},
						{Note: &symbols.Note{
							Embellishment: &symbols.Embellishment{Type: symbols.Doubling},
						}},
						{Note: &symbols.Note{Pitch: common.LowA, Length: common.Eighth}},
					},
				}
				f.wantErr = true
			},
		},
		{
			name: "all preceding symbols are notes",
			prepare: func(f *fields) {
				f.measure = &music_model.Measure{
					Time: nil,
					Symbols: []*music_model.Symbol{
						{Note: &symbols.Note{Pitch: common.LowA, Length: common.Eighth}},
						{Note: &symbols.Note{Pitch: common.LowA, Length: common.Eighth}},
						{Note: &symbols.Note{Pitch: common.LowA, Length: common.Eighth}},
					},
				}
				f.wantErr = false
			},
			after: func(f *fields) {
				g.Expect(f.measure.Symbols).To(HaveLen(5))
				g.Expect(f.measure.Symbols[0].Tuplet).To(Equal(&tuplet.Tuplet{
					BoundaryType: tuplet.Start,
					VisibleNotes: 3,
					PlayedNotes:  2,
				}))
				g.Expect(f.measure.Symbols[4].Tuplet).To(Equal(&tuplet.Tuplet{
					BoundaryType: tuplet.End,
					VisibleNotes: 3,
					PlayedNotes:  2,
				}))
			},
		},
		{
			name: "if there is already a tuplet start, don't add another one",
			prepare: func(f *fields) {
				f.measure = &music_model.Measure{
					Time: nil,
					Symbols: []*music_model.Symbol{
						{
							Tuplet: &tuplet.Tuplet{
								BoundaryType: tuplet.Start,
								VisibleNotes: 3,
								PlayedNotes:  2,
							},
						},
						{Note: &symbols.Note{Pitch: common.LowA, Length: common.Eighth}},
						{Note: &symbols.Note{Pitch: common.LowA, Length: common.Eighth}},
						{Note: &symbols.Note{Pitch: common.LowA, Length: common.Eighth}},
					},
				}
				f.wantErr = false
			},
			after: func(f *fields) {
				g.Expect(f.measure.Symbols).To(HaveLen(5))
				g.Expect(f.measure.Symbols[0].Tuplet).To(Equal(&tuplet.Tuplet{
					BoundaryType: tuplet.Start,
					VisibleNotes: 3,
					PlayedNotes:  2,
				}))
				g.Expect(f.measure.Symbols[4].Tuplet).To(Equal(&tuplet.Tuplet{
					BoundaryType: tuplet.End,
					VisibleNotes: 3,
					PlayedNotes:  2,
				}))
			},
		},
		{
			name: "if there is a tuplet end, a start mus be added",
			prepare: func(f *fields) {
				f.measure = &music_model.Measure{
					Time: nil,
					Symbols: []*music_model.Symbol{
						{
							Tuplet: &tuplet.Tuplet{
								BoundaryType: tuplet.End,
								VisibleNotes: 7,
								PlayedNotes:  6,
							},
						},
						{Note: &symbols.Note{Pitch: common.LowA, Length: common.Eighth}},
						{Note: &symbols.Note{Pitch: common.LowA, Length: common.Eighth}},
						{Note: &symbols.Note{Pitch: common.LowA, Length: common.Eighth}},
					},
				}
				f.wantErr = false
			},
			after: func(f *fields) {
				g.Expect(f.measure.Symbols).To(HaveLen(6))
				g.Expect(f.measure.Symbols[0].Tuplet).To(Equal(&tuplet.Tuplet{
					BoundaryType: tuplet.End,
					VisibleNotes: 7,
					PlayedNotes:  6,
				}))
				g.Expect(f.measure.Symbols[1].Tuplet).To(Equal(&tuplet.Tuplet{
					BoundaryType: tuplet.Start,
					VisibleNotes: 3,
					PlayedNotes:  2,
				}))
				g.Expect(f.measure.Symbols[5].Tuplet).To(Equal(&tuplet.Tuplet{
					BoundaryType: tuplet.End,
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
