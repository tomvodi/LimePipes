package model

import (
	"time"
)

type ImportFile struct {
	ID           uint64 `gorm:"primaryKey"`
	Name         string
	OriginalPath string
	Hash         string
	Data         []byte // file content from original file
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
