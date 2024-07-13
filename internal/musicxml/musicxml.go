package musicxml

import (
	"encoding/xml"
	"github.com/tomvodi/limepipes/internal/musicxml/model"
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
