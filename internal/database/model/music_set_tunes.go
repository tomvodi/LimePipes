package model

type MusicSetTunes struct {
	MusicSetID uint64 `gorm:"primaryKey"`
	TuneID     uint64 `gorm:"primaryKey"`
	Order      uint   `gorm:"not null"`
}
