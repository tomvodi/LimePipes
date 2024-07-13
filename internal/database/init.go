package database

import (
	"github.com/tomvodi/limepipes/internal/database/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetInitSqliteDb(name string) (*gorm.DB, error) {
	dialect := sqlite.Open(name)
	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if res := db.Exec("PRAGMA foreign_keys = ON", nil); res.Error != nil {
		return nil, res.Error
	}

	if err = db.AutoMigrate(
		&model.MusicSet{},
		&model.Tune{},
		&model.MusicSetTunes{},
		&model.TuneFile{},
		&model.ImportFile{},
		&model.TuneType{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
