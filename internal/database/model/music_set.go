package model

type MusicSet struct {
	BaseModel
	Title        string
	Description  string
	Creator      string
	Tunes        []Tune `gorm:"many2many:music_set_tunes;constraint:OnUpdate:CASCADE;OnDelete:RESTRICT"`
	ImportFileId int64
}
