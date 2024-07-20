package model

import "github.com/SamuelTissot/sqltime"

type BaseModel struct {
	ID        int64        `gorm:"primary_key" copier:"Id"`
	CreatedAt sqltime.Time `gorm:"type:timestamp"`
	UpdatedAt sqltime.Time `gorm:"type:timestamp"`
}
