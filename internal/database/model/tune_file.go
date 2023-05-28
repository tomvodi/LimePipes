package model

import (
	"banduslib/internal/common/music_model"
	"banduslib/internal/database/model/file_type"
	"bytes"
	"encoding/gob"
	"fmt"
	"gorm.io/gorm"
)

type TuneFile struct {
	gorm.Model
	TuneID uint
	Type   file_type.Type
	Data   []byte
}

func (t *TuneFile) MusicModel() (music_model.MusicModel, error) {
	if t.Type != file_type.MusicModel {
		return nil, fmt.Errorf("tune file has type %s not type %s",
			t.Type.String(), file_type.MusicModel.String(),
		)
	}

	if t.Data == nil {
		return nil, fmt.Errorf("can't get music model from tune file as data is empty")
	}

	buf := bytes.NewBuffer(t.Data)
	dec := gob.NewDecoder(buf)

	var mumo music_model.MusicModel

	if err := dec.Decode(&mumo); err != nil {
		return nil, err
	}

	return mumo, nil
}

func TuneFileFromMusicModel(muMod music_model.MusicModel) (*TuneFile, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(muMod); err != nil {
		return nil, err
	}

	tf := &TuneFile{
		Type: file_type.MusicModel,
		Data: buf.Bytes(),
	}

	return tf, nil
}
