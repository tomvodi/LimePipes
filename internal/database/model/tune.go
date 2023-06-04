package model

import (
	"time"
)

type Tune struct {
	ID           uint64 `gorm:"primaryKey"`
	Title        string
	TuneTypeId   *uint64
	TuneType     *TuneType
	TimeSig      string
	Composer     string
	Arranger     string
	Sets         []MusicSet  `gorm:"many2many:music_set_tunes;constraint:OnUpdate:CASCADE;"`
	Files        []*TuneFile `gorm:"constraint:OnDelete:CASCADE;"`
	ImportFileId uint64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
