package model

type MusicSetTunes struct {
	BaseModel
	MusicSetID int64
	MusicSet   MusicSet
	TuneID     int64
	Tune       Tune
	Order      uint `gorm:"not null"`
}
