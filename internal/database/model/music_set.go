package model

import (
	"time"
)

type MusicSet struct {
	ID           int64 `gorm:"primaryKey"`
	Title        string
	Description  string
	Creator      string
	Tunes        []Tune `gorm:"many2many:music_set_tunes;constraint:OnUpdate:CASCADE;OnDelete:RESTRICT"`
	ImportFileId int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
