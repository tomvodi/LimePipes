package helper

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/utils"
	. "github.com/onsi/gomega"
	"testing"
)

func Test_fixComposer(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		composer     string
		wantComposer string
		wantArranger string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: "arr.",
			prepare: func(f *fields) {
				f.composer = "arr. Willi World"
				f.wantComposer = ""
				f.wantArranger = "Willi World"
			},
		},
		{
			name: "Arr.",
			prepare: func(f *fields) {
				f.composer = "Arr. Willi World"
				f.wantComposer = ""
				f.wantArranger = "Willi World"
			},
		},
		{
			name: "arr:",
			prepare: func(f *fields) {
				f.composer = "arr: Willi World"
				f.wantComposer = ""
				f.wantArranger = "Willi World"
			},
		},
		{
			name: "Arr:",
			prepare: func(f *fields) {
				f.composer = "Arr: Willi World"
				f.wantComposer = ""
				f.wantArranger = "Willi World"
			},
		},
		{
			name: "arr/",
			prepare: func(f *fields) {
				f.composer = "arr/ Willi World"
				f.wantComposer = ""
				f.wantArranger = "Willi World"
			},
		},
		{
			name: "Arr/",
			prepare: func(f *fields) {
				f.composer = "Arr/ Willi World"
				f.wantComposer = ""
				f.wantArranger = "Willi World"
			},
		},
		{
			name: "arr",
			prepare: func(f *fields) {
				f.composer = "arr Willi World"
				f.wantComposer = ""
				f.wantArranger = "Willi World"
			},
		},
		{
			name: "Arr",
			prepare: func(f *fields) {
				f.composer = "Arr Willi World"
				f.wantComposer = ""
				f.wantArranger = "Willi World"
			},
		},
		{
			name: "arranged by",
			prepare: func(f *fields) {
				f.composer = "arranged by Willi World"
				f.wantComposer = ""
				f.wantArranger = "Willi World"
			},
		},
		{
			name: "tuneTitle with whitespace in beginning or end",
			prepare: func(f *fields) {
				f.composer = "  arranged by Willi World  "
				f.wantComposer = ""
				f.wantArranger = "Willi World"
			},
		},
		{
			name: "Arranged by",
			prepare: func(f *fields) {
				f.composer = "Arranged by Willi World"
				f.wantComposer = ""
				f.wantArranger = "Willi World"
			},
		},
		{
			name: "arrangement",
			prepare: func(f *fields) {
				f.composer = "arrangement Willi World"
				f.wantComposer = ""
				f.wantArranger = "Willi World"
			},
		},
		{
			name: "Arrangement",
			prepare: func(f *fields) {
				f.composer = "Arrangement Willi World"
				f.wantComposer = ""
				f.wantArranger = "Willi World"
			},
		},
		{
			name: "J. A. Barrie",
			prepare: func(f *fields) {
				f.composer = "J. A. Barrie"
				f.wantComposer = "J. A. Barrie"
				f.wantArranger = ""
			},
		},
		{
			name: "Arr. P/M A. MacDonald, Scots Guards",
			prepare: func(f *fields) {
				f.composer = "Arr. P/M A. MacDonald, Scots Guards"
				f.wantComposer = "Scots Guards"
				f.wantArranger = "P/M A. MacDonald"
			},
		},
		{
			name: "R. Williamson - arr. P/M J.B. Salad",
			prepare: func(f *fields) {
				f.composer = "R. Williamson - arr. P/M J.B. Salad"
				f.wantComposer = "R. Williamson"
				f.wantArranger = "P/M J.B. Salad"
			},
		},
		{
			name: "R. Williamson, arr. P/M J.B. Salad",
			prepare: func(f *fields) {
				f.composer = "R. Williamson, arr. P/M J.B. Salad"
				f.wantComposer = "R. Williamson"
				f.wantArranger = "P/M J.B. Salad"
			},
		},
		{
			name: "Composer Unknown [Arr: Pipe Major R Hear]",
			prepare: func(f *fields) {
				f.composer = "Composer Unknown [Arr: Pipe Major R Hear]"
				f.wantComposer = "Composer Unknown"
				f.wantArranger = "Pipe Major R Hear"
			},
		},
		{
			name: "Traditional (arr: College of PipingTutorVol 1)",
			prepare: func(f *fields) {
				f.composer = "Traditional (arr: College of PipingTutorVol 1)"
				f.wantComposer = "Traditional"
				f.wantArranger = "College of PipingTutorVol 1"
			},
		},
		{
			name: "Arr. Jason Monday, David Pear",
			prepare: func(f *fields) {
				f.composer = "Arr. Jason Monday, David Pear"
				f.wantComposer = "David Pear"
				f.wantArranger = "Jason Monday"
			},
		},
		{
			name: "Trad. arr. E Mule",
			prepare: func(f *fields) {
				f.composer = "Trad. arr. E Mule"
				f.wantComposer = "Trad"
				f.wantArranger = "E Mule"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			tune := &music_model.Tune{
				Composer: f.composer,
			}
			fixComposerArranger(tune)
			g.Expect(tune.Composer).To(Equal(f.wantComposer))
			g.Expect(tune.Arranger).To(Equal(f.wantArranger))
		})
	}
}

func Test_fixComposerTrad(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		composer     string
		wantComposer string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: "Trad.",
			prepare: func(f *fields) {
				f.composer = "Trad."
				f.wantComposer = "Traditional"
			},
		},
		{
			name: "Trad",
			prepare: func(f *fields) {
				f.composer = "Trad"
				f.wantComposer = "Traditional"
			},
		},
		{
			name: "trad",
			prepare: func(f *fields) {
				f.composer = "trad"
				f.wantComposer = "Traditional"
			},
		},
		{
			name: " trad ",
			prepare: func(f *fields) {
				f.composer = " trad "
				f.wantComposer = "Traditional"
			},
		},
		{
			name: "Tradontio",
			prepare: func(f *fields) {
				f.composer = "Tradontio"
				f.wantComposer = "Tradontio"
			},
		},
		{
			name: "Is Trad. Tune",
			prepare: func(f *fields) {
				f.composer = "Is Trad. Tune"
				f.wantComposer = "Is Trad. Tune"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			tune := &music_model.Tune{
				Composer: f.composer,
			}
			fixComposerTrad(tune)
			g.Expect(tune.Composer).To(Equal(f.wantComposer))
		})
	}
}

func Test_removeTimeSigFromTuneType(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		tuneType     string
		wantTuneType string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: " 2/4 Hornpipe ",
			prepare: func(f *fields) {
				f.tuneType = " 2/4 Hornpipe "
				f.wantTuneType = "Hornpipe"
			},
		},
		{
			name: " Hornpipe 2/4 ",
			prepare: func(f *fields) {
				f.tuneType = " Hornpipe 2/4 "
				f.wantTuneType = "Hornpipe"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			tune := &music_model.Tune{
				Type: f.tuneType,
			}
			removeTimeSigFromTuneType(tune)
			g.Expect(tune.Type).To(Equal(f.wantTuneType))
		})
	}
}

func Test_capitalizeTuneType(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		tuneType     string
		wantTuneType string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: " reel ",
			prepare: func(f *fields) {
				f.tuneType = " reel"
				f.wantTuneType = "Reel"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			tune := &music_model.Tune{
				Type: f.tuneType,
			}
			capitalizeTuneType(tune)
			g.Expect(tune.Type).To(Equal(f.wantTuneType))
		})
	}
}

func Test_trimSpaces(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		tuneType         string
		tuneTitle        string
		tuneComposer     string
		tuneArranger     string
		wantTuneType     string
		wantTuneTitle    string
		wantTuneComposer string
		wantTuneArranger string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: " trim all ",
			prepare: func(f *fields) {
				f.tuneType = " type "
				f.tuneTitle = " title "
				f.tuneComposer = " composer "
				f.tuneArranger = " arranger"
				f.wantTuneType = "type"
				f.wantTuneTitle = "title"
				f.wantTuneComposer = "composer"
				f.wantTuneArranger = "arranger"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			tune := &music_model.Tune{
				Type:     f.tuneType,
				Title:    f.tuneTitle,
				Composer: f.tuneComposer,
				Arranger: f.tuneArranger,
			}
			trimSpaces(tune)
			g.Expect(tune.Title).To(Equal(f.wantTuneTitle))
			g.Expect(tune.Type).To(Equal(f.wantTuneType))
			g.Expect(tune.Composer).To(Equal(f.wantTuneComposer))
			g.Expect(tune.Arranger).To(Equal(f.wantTuneArranger))
		})
	}
}

func Test_fixTitle(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		tuneTitle     string
		wantTuneTitle string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: "the_tune_title",
			prepare: func(f *fields) {
				f.tuneTitle = "the_tune_title"
				f.wantTuneTitle = "The Tune Title"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			tune := &music_model.Tune{
				Title: f.tuneTitle,
			}
			fixTitle(tune)
			g.Expect(tune.Title).To(Equal(f.wantTuneTitle))
		})
	}
}
