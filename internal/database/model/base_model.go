package model

import (
	"github.com/SamuelTissot/sqltime"
	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid()" copier:"ID"`
	CreatedAt sqltime.Time `gorm:"type:timestamp"`
	UpdatedAt sqltime.Time `gorm:"type:timestamp"`
}
