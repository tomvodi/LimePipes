package model

import (
	"time"
)

type MusicSet struct {
	ID          uint64 `gorm:"primaryKey"`
	Title       string `gorm:"unique;not null;default:null"`
	Description string
	Creator     string
	Tunes       []Tune `gorm:"many2many:music_set_tunes;constraint:OnUpdate:CASCADE;OnDelete:RESTRICT"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
