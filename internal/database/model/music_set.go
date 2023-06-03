package model

import (
	"time"
)

type MusicSet struct {
	ID           uint64 `gorm:"primaryKey"`
	Title        string
	Description  string
	Creator      string
	Tunes        []Tune `gorm:"many2many:music_set_tunes;constraint:OnUpdate:CASCADE;OnDelete:RESTRICT"`
	ImportFileId uint64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
