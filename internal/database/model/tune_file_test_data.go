package model

import (
	"github.com/tomvodi/limepipes/internal/common"
	"github.com/tomvodi/limepipes/internal/common/music_model"
	"github.com/tomvodi/limepipes/internal/common/music_model/barline"
	"github.com/tomvodi/limepipes/internal/common/music_model/import_message"
	"github.com/tomvodi/limepipes/internal/common/music_model/symbols"
	"github.com/tomvodi/limepipes/internal/common/music_model/symbols/accidental"
	"github.com/tomvodi/limepipes/internal/common/music_model/symbols/embellishment"
	"github.com/tomvodi/limepipes/internal/common/music_model/symbols/tie"
)

func TestMusicModelTune(title string) *music_model.Tune {
	return &music_model.Tune{
		Title:      title,
		Type:       "march",
		Composer:   "someone",
		Arranger:   "someone arranged it",
		Footer:     []string{"footer tune 1"},
		Comments:   []string{"comment 1", "comment 2"},
		InLineText: []string{"inline text 1", "inline text 2"},
		Tempo:      80,
		Measures: []*music_model.Measure{
			{
				LeftBarline: &barline.Barline{
					Type:     barline.Heavy,
					Timeline: barline.Segno,
				},
				RightBarline: &barline.Barline{
					Type:     barline.Heavy,
					Timeline: barline.DacapoAlFine,
				},
				Time: &music_model.TimeSignature{
					Beats:    2,
					BeatType: 4,
				},
				Symbols: []*music_model.Symbol{
					{
						Note: &symbols.Note{
							Pitch:      common.LowA,
							Length:     common.Quarter,
							Dots:       2,
							Accidental: accidental.Natural,
							Fermata:    true,
							Tie:        tie.Start,
							Embellishment: &embellishment.Embellishment{
								Type:    embellishment.Doubling,
								Pitch:   common.E,
								Variant: embellishment.Half,
								Weight:  embellishment.Light,
							},
							Movement: nil,
							Comment:  "",
						},
						Rest:        nil,
						Tuplet:      nil,
						TimeLine:    nil,
						TempoChange: 0,
					},
				},
				Comments:   []string{"comment measure 1", "comment measure 2"},
				InLineText: []string{"inline text measure 1", "inline text measure 2"},
				ImportMessages: []*import_message.ImportMessage{
					{
						Symbol: "^te",
						Type:   import_message.Warning,
						Text:   "some warning",
						Fix:    import_message.SkipSymbol,
					},
				},
			},
		},
	}
}
