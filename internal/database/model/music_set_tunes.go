package model

type MusicSetTunes struct {
	MusicSetID uint64 `gorm:"primaryKey"`
	MusicSet   MusicSet
	TuneID     uint64 `gorm:"primaryKey"`
	Tune       Tune
	Order      uint `gorm:"not null"`
}
