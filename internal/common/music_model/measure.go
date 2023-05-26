package music_model

import "banduslib/internal/common/music_model/barline"

type Measure struct {
	LeftBarline  *barline.Barline `yaml:"leftBarline,omitempty"`
	RightBarline *barline.Barline `yaml:"rightBarline,omitempty"`
	Time         *TimeSignature   `yaml:"time,omitempty"`
	Symbols      []*Symbol        `yaml:"symbols,omitempty"`
	Comments     []string         `yaml:"comments,omitempty"`
	InLineText   []string         `yaml:"inLineText,omitempty"`
}
