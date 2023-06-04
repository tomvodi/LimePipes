package musicxml

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/musicxml/model"
	"encoding/xml"
	"io"
)

func WriteScore(score *model.Score, writer io.Writer) error {
	data, err := xml.MarshalIndent(score, " ", "  ")
	if err != nil {
		return err
	}

	data = append([]byte(musicXMLHeader), data...)
	if _, err := writer.Write(data); err != nil {
		return err
	}

	return nil
}

func ReadScore(reader io.Reader) (*model.Score, error) {
	fileData, _ := io.ReadAll(reader)

	score := &model.Score{}
	err := xml.Unmarshal(fileData, score)
	if err != nil {
		return nil, err
	}

	return score, nil
}

func ScoreFromMusicModelTune(tune *music_model.Tune) (*model.Score, error) {
	measures := []model.Measure{}
	for i, measure := range tune.Measures {
		xmlMeasure := xmlMeasureFromMusicModelMeasure(measure, i)
		measures = append(measures, xmlMeasure)
	}

	score := &model.Score{
		XMLName: xml.Name{
			Local: "score-partwise",
		},
		Version: "3.1",
		PartList: model.ScorePartList{
			XMLName: xml.Name{
				Local: "part-list",
			},
			Parts: []model.ScorePart{
				{
					XMLName: xml.Name{
						Local: "score-part",
					},
					Id:   "P1",
					Name: "Bagpipe",
					Instrument: model.ScoreInstrument{
						XMLName: xml.Name{
							Local: "score-instrument",
						},
						Id:   "P1-I1",
						Name: "Bagpipe",
					},
					MidiDevice: model.MidiDevice{
						XMLName: xml.Name{
							Local: "midi-device",
						},
						Id:   "P1-I1",
						Port: 1,
					},
					MidiInstrument: model.MidiInstrument{
						XMLName: xml.Name{
							Local: "midi-instrument",
						},
						Id:      "P1-I1",
						Channel: 1,
						Program: 110,
						Volume:  78.7402,
						Pan:     0,
					},
				},
			},
		},
		Part: model.Part{
			XMLName: xml.Name{
				Local: "part",
			},
			Id:       "P1",
			Measures: measures,
		},
	}

	return score, nil
}

func xmlMeasureFromMusicModelMeasure(measure *music_model.Measure, id int) model.Measure {
	xmlMeasure := model.Measure{
		XMLName: xml.Name{
			Local: "measure",
		},
		Number: id,
	}
	return xmlMeasure
}
