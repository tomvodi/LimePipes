package model

import (
	"time"
)

type ImportFile struct {
	ID           int64 `gorm:"primaryKey" copier:"Id"`
	Name         string
	OriginalPath string
	Hash         string
	Data         []byte // file content from original file
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
