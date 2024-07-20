package model

type TuneType struct {
	BaseModel
	Name string `gorm:"unique"`
}
