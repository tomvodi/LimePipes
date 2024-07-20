package database

import (
	"fmt"
	"github.com/tomvodi/limepipes/internal/config"
	"github.com/tomvodi/limepipes/internal/database/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetInitPostgreSQLDB(
	dbConf config.DbConfig,
) (*gorm.DB, error) {
	dsnTpl := "postgres://%s:%s@%s:%s/%s?sslmode=%s"
	dsn := fmt.Sprintf(dsnTpl,
		dbConf.User,
		dbConf.Password,
		dbConf.Host,
		dbConf.Port,
		dbConf.DbName,
		dbConf.SslMode,
	)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
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
