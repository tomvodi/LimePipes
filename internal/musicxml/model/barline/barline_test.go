package barline

import (
	"banduslib/internal/common/music_model/barline"
	"banduslib/internal/utils"
	. "github.com/onsi/gomega"
	"testing"
)

func Test_convertBarlineType(t *testing.T) {
	utils.SetupConsoleLogger()
	g := NewGomegaWithT(t)
	type fields struct {
		barlineType barline.BarlineType
		want        Style
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
	}{
		{
			name: "Regular",
			prepare: func(f *fields) {
				f.barlineType = barline.Regular
				f.want = Regular
			},
		},
		{
			name: "Heavy",
			prepare: func(f *fields) {
				f.barlineType = barline.Heavy
				f.want = Heavy
			},
		},
		{
			name: "HeavyHeavy",
			prepare: func(f *fields) {
				f.barlineType = barline.HeavyHeavy
				f.want = HeavyHeavy
			},
		},
		{
			name: "LightHeavy",
			prepare: func(f *fields) {
				f.barlineType = barline.LightHeavy
				f.want = LightHeavy
			},
		},
		{
			name: "HeavyLight",
			prepare: func(f *fields) {
				f.barlineType = barline.HeavyLight
				f.want = HeavyLight
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := &fields{}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			got := convertBarlineType(f.barlineType)
			g.Expect(got).To(Equal(f.want))
		})
	}

}
