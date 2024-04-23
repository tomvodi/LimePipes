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

func (m *Measure) HasOnlyAttributes() bool {
	onlyAttributes := true
	for _, symbol := range m.Symbols {
		if symbol.IsNote() && symbol.IsValidNote() {
			onlyAttributes = false
		}
		if symbol.Rest != nil {
			onlyAttributes = false
		}
	}
	return onlyAttributes
}

func (m *Measure) HasSymbols() bool {
	return len(m.Symbols) > 0
}

// ContainsTimelineSymbols return true if measure contains at least one timeline symbol
func (m *Measure) ContainsTimelineSymbols() bool {
	for _, sym := range m.Symbols {
		if sym.IsTimeline() {
			return true
		}
	}

	return false
}

// IsOneTimeline returns true, if measure begins and ens with timeline and
// there are no timeline symbol in between
func (m *Measure) IsOneTimeline() bool {
	if len(m.Symbols) < 2 {
		return false
	}

	timelineSymCnt := 0
	for _, symbol := range m.Symbols {
		if symbol.IsTimeline() {
			timelineSymCnt++
		}
	}
	if timelineSymCnt != 2 {
		return false
	}

	firstSym := m.Symbols[0]
	lastSym := m.Symbols[len(m.Symbols)-1]
	if firstSym.IsTimelineStart() &&
		lastSym.IsTimelineEnd() {
		return true
	}

	return false
}
