package musicxml

import (
	"github.com/tomvodi/limepipes/internal/musicxml/model"
	"os"
	"testing"
)

func TestWriteScore(t *testing.T) {
	score := &model.Score{
		PartList: model.ScorePartList{
			Parts: []model.ScorePart{
				{
					ID:   "P1",
					Name: "Bagpipe",
					Instrument: model.ScoreInstrument{
						ID:   "P1-I1",
						Name: "Bagpipe",
					},
					MidiDevice: model.MidiDevice{
						ID:   "P1-I1",
						Port: 1,
					},
					MidiInstrument: model.MidiInstrument{
						ID:       "P1-I1",
						Channel:  1,
						Programm: 110,
						Volume:   78.7402,
						Pan:      0,
					},
				},
			},
		},
		Credits: []model.Credit{
			{
				Page: 1,
				Type: model.CreditTypeComposer,
				Words: model.CreditWords{
					Value: "Score composer",
				},
			},
			{
				Page: 1,
				Type: model.CreditTypeTitle,
				Words: model.CreditWords{
					Value: "My Favorite Tune",
				},
			},
		},
	}
	err := WriteScore(score, os.Stdout)
	if err != nil {
		t.Errorf("WriteScore() error = %v", err)
		return
	}

}
