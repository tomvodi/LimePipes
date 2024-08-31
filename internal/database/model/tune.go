package model

import "github.com/google/uuid"

type Tune struct {
	BaseModel
	Title        string
	TuneTypeID   *uuid.UUID
	TuneType     *TuneType
	TimeSig      string
	Composer     string
	Arranger     string
	Sets         []MusicSet  `gorm:"many2many:music_set_tunes;constraint:OnUpdate:CASCADE;"`
	Files        []*TuneFile `gorm:"constraint:OnDelete:CASCADE;"`
	ImportFileID uuid.UUID
}
