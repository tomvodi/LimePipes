package model

import (
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/barline"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/length"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/measure"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/pitch"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/symbols/accidental"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/symbols/embellishment"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/symbols/tie"
	"github.com/tomvodi/limepipes-music-model/musicmodel/v1/tune"
)

func TestMusicModelTune(title string) *tune.Tune {
	return &tune.Tune{
		Title:      title,
		Type:       "march",
		Composer:   "someone",
		Arranger:   "someone arranged it",
		Footer:     []string{"footer tune 1"},
		Comments:   []string{"comment 1", "comment 2"},
		InLineText: []string{"inline text 1", "inline text 2"},
		Tempo:      80,
		Measures: []*measure.Measure{
			{
				LeftBarline: &barline.Barline{
					Type: barline.Type_Heavy,
					Time: barline.Time_Segno,
				},
				RightBarline: &barline.Barline{
					Type: barline.Type_Heavy,
					Time: barline.Time_DacapoAlFine,
				},
				Time: &measure.TimeSignature{
					Beats:    2,
					BeatType: 4,
				},
				Symbols: []*symbols.Symbol{
					{
						Note: &symbols.Note{
							Pitch:      pitch.Pitch_LowA,
							Length:     length.Length_Quarter,
							Dots:       2,
							Accidental: accidental.Accidental_Natural,
							Fermata:    true,
							Tie:        tie.Tie_Start,
							Embellishment: &embellishment.Embellishment{
								Type:    embellishment.Type_Doubling,
								Pitch:   pitch.Pitch_E,
								Variant: embellishment.Variant_Half,
								Weight:  embellishment.Weight_Light,
							},
							Movement: nil,
							Comment:  "",
						},
						Rest:        nil,
						Tuplet:      nil,
						Timeline:    nil,
						TempoChange: nil,
					},
				},
				Comments:   []string{"comment measure 1", "comment measure 2"},
				InlineText: []string{"inline text measure 1", "inline text measure 2"},
				ImportMessages: []*measure.ImportMessage{
					{
						Symbol:   "^te",
						Severity: measure.Severity_Warning,
						Text:     "some warning",
						Fix:      measure.Fix_SkipSymbol,
					},
				},
			},
		},
	}
}
