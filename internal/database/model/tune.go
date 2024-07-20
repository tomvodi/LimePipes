package model

type Tune struct {
	BaseModel
	Title        string
	TuneTypeId   *int64
	TuneType     *TuneType
	TimeSig      string
	Composer     string
	Arranger     string
	Sets         []MusicSet  `gorm:"many2many:music_set_tunes;constraint:OnUpdate:CASCADE;"`
	Files        []*TuneFile `gorm:"constraint:OnDelete:CASCADE;"`
	ImportFileId int64
}
