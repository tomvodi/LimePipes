package music_model

import "banduslib/internal/common/music_model/import_message"

type Tune struct {
	Title      string     `yaml:"title"`
	Type       string     `yaml:"type,omitempty"`
	Composer   string     `yaml:"composer,omitempty"`
	Footer     []string   `yaml:"footer,omitempty"`
	Comments   []string   `yaml:"comments,omitempty"`
	InLineText []string   `yaml:"inLineText,omitempty"`
	Tempo      uint64     `yaml:"tempo"`
	Measures   []*Measure `yaml:"measures"`
}

func (t *Tune) ImportMessages() []*import_message.ImportMessage {
	var messages []*import_message.ImportMessage

	for _, measure := range t.Measures {
		for _, message := range measure.ImportMessages {
			messages = append(messages, message)
		}
	}

	return messages
}

func (t *Tune) FirstTimeSignature() *TimeSignature {
	for _, measure := range t.Measures {
		if measure.Time != nil {
			return measure.Time
		}
	}

	return nil
}
