package model

import "time"

type TuneType struct {
	ID        uint64 `gorm:"primaryKey"`
	Name      string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
