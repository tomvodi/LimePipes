package model

import (
	"time"
)

type Tune struct {
	ID        uint64 `gorm:"primaryKey"`
	Title     string `gorm:"unique;not null;default:null"`
	Type      string
	TimeSig   string
	Composer  string
	Arranger  string
	Sets      []MusicSet ` gorm:"many2many:music_set_tunes;constraint:OnUpdate:CASCADE;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
