package model

type MusicSetTunes struct {
	ID         uint64 `gorm:"primaryKey"`
	MusicSetID uint64
	MusicSet   MusicSet
	TuneID     uint64
	Tune       Tune
	Order      uint `gorm:"not null"`
}
