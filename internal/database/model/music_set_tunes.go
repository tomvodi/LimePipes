package model

import "github.com/google/uuid"

type MusicSetTunes struct {
	BaseModel
	MusicSetID uuid.UUID
	MusicSet   MusicSet
	TuneID     uuid.UUID
	Tune       Tune
	Order      uint `gorm:"not null"`
}
