package database

import (
	. "github.com/onsi/gomega"
	"github.com/tomvodi/limepipes/internal/apigen/apimodel"
	"github.com/tomvodi/limepipes/internal/utils"
	"testing"
)

func Test_musicSetTitleFromTunes(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		tunes []*apimodel.ImportTune
		want  string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: "tune types and time signatures are the same",
			prepare: func(f *fields) {
				f.tunes = []*apimodel.ImportTune{
					{
						Type:    "March",
						TimeSig: "6/8",
					},
					{
						Type:    "March",
						TimeSig: "6/8",
					},
				}
				f.want = "6/8 March Set"
			},
		},
		{
			name: "tune types are the same",
			prepare: func(f *fields) {
				f.tunes = []*apimodel.ImportTune{
					{
						Type:    "March",
						TimeSig: "6/8",
					},
					{
						Type:    "March",
						TimeSig: "4/4",
					},
				}
				f.want = "March Set"
			},
		},
		{
			name: "MSR set",
			prepare: func(f *fields) {
				f.tunes = []*apimodel.ImportTune{
					{
						Type:    "March",
						TimeSig: "6/8",
					},
					{
						Type:    "Strathspey",
						TimeSig: "4/4",
					},
					{
						Type:    "Reel",
						TimeSig: "4/4",
					},
				}
				f.want = "MSR Set"
			},
		},
		{
			name: "Different tune types",
			prepare: func(f *fields) {
				f.tunes = []*apimodel.ImportTune{
					{
						Type: "Slow Air",
					},
					{
						Type: "Hornpipe",
					},
					{
						Type: "Jig",
					},
				}
				f.want = "Slow Air - Hornpipe - Jig"
			},
		},
		{
			name: "Different and missing tune types",
			prepare: func(f *fields) {
				f.tunes = []*apimodel.ImportTune{
					{
						Type: "Slow Air",
					},
					{
						Type: "",
					},
					{
						Type: "Jig",
					},
				}
				f.want = "Slow Air - Unknown Type - Jig"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			setName := musicSetTitleFromTunes(f.tunes)
			g.Expect(setName).To(Equal(f.want))
		})
	}
}
