package music_model

import (
	"banduslib/internal/common/music_model/barline"
	"banduslib/internal/common/music_model/import_message"
	"reflect"
)

type Measure struct {
	LeftBarline    *barline.Barline                `yaml:"leftBarline,omitempty"`
	RightBarline   *barline.Barline                `yaml:"rightBarline,omitempty"`
	Time           *TimeSignature                  `yaml:"time,omitempty"`
	Symbols        []*Symbol                       `yaml:"symbols,omitempty"`
	Comments       []string                        `yaml:"comments,omitempty"`
	InLineText     []string                        `yaml:"inLineText,omitempty"`
	ImportMessages []*import_message.ImportMessage `yaml:"-"`
}

func (m *Measure) AddMessage(msg *import_message.ImportMessage) {
	if m.ImportMessages == nil {
		m.ImportMessages = []*import_message.ImportMessage{}
	}

	m.ImportMessages = append(m.ImportMessages, msg)
}

func (m *Measure) IsNil() bool {
	return reflect.DeepEqual(m, &Measure{})
}
