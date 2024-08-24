package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/google/uuid"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/tune"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/file_type"
)

type TuneFile struct {
	TuneID uuid.UUID      `gorm:"primaryKey"`
	Type   file_type.Type `gorm:"primaryKey"`

	// SingleTuneData is true if the data is for a single tune and not the whole file
	// e.g. for a music model tune from a .bww file which supports multiple tunes in one file,
	// there can be a tune file data for the whole file and one for the single tune.
	SingleTuneData bool `gorm:"primaryKey"`

	Data []byte
}

func (t *TuneFile) MusicModelTune() (*tune.Tune, error) {
	if t.Type != file_type.Type_MUSIC_MODEL {
		return nil, fmt.Errorf("tune file has type %s not type %s",
			t.Type.String(), file_type.Type_MUSIC_MODEL.String(),
		)
	}

	if t.Data == nil {
		return nil, fmt.Errorf("can't get music model tune from tune file as data is empty")
	}

	buf := bytes.NewBuffer(t.Data)
	dec := gob.NewDecoder(buf)

	tn := &tune.Tune{}

	if err := dec.Decode(tn); err != nil {
		return nil, err
	}

	return tn, nil
}

func TuneFileFromMusicModelTune(tune *tune.Tune) (*TuneFile, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(tune); err != nil {
		return nil, err
	}

	tf := &TuneFile{
		Type: file_type.Type_MUSIC_MODEL,
		Data: buf.Bytes(),
	}

	return tf, nil
}
