package model

type MusicSetTunes struct {
	ID         int64 `gorm:"primaryKey"`
	MusicSetID int64
	MusicSet   MusicSet
	TuneID     int64
	Tune       Tune
	Order      uint `gorm:"not null"`
}
