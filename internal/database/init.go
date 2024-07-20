package database

import (
	"fmt"
	"github.com/tomvodi/limepipes/internal/config"
	"github.com/tomvodi/limepipes/internal/database/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func migrateDb(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.MusicSet{},
		&model.Tune{},
		&model.MusicSetTunes{},
		&model.TuneFile{},
		&model.ImportFile{},
		&model.TuneType{},
	)
}

func GetInitPostgreSQLDB(
	dbConf config.DbConfig,
) (*gorm.DB, error) {
	dsnTpl := "host=%s port=%s dbname=%s user=%s password=%s sslmode=%s TimeZone=%s"
	dsn := fmt.Sprintf(dsnTpl,
		dbConf.Host,
		dbConf.Port,
		dbConf.DbName,
		dbConf.User,
		dbConf.Password,
		dbConf.SslMode,
		dbConf.TimeZone,
	)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = migrateDb(db); err != nil {
		return nil, err
	}

	return db, nil
}

func GetInitTestPostgreSQLDB(
	dbConf config.DbConfig,
	testDbName string,
) (*gorm.DB, error) {
	// Connect without database name to create test database
	dsnTpl := "host=%s port=%s user=%s password=%s sslmode=%s TimeZone=%s"
	dsn := fmt.Sprintf(dsnTpl,
		dbConf.Host,
		dbConf.Port,
		dbConf.User,
		dbConf.Password,
		dbConf.SslMode,
		dbConf.TimeZone,
	)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Disconnect all other sessions to be able to drop the database
	disconnectCmd := fmt.Sprintf(`SELECT 
	pg_terminate_backend(pid) 
	FROM 
	pg_stat_activity
	WHERE
	-- don't kill my own connection!
	pid <> pg_backend_pid()
	-- don't kill the connections to other databases
	AND datname = '%s'
	;`, testDbName)
	db.Exec(disconnectCmd)

	dropDatabaseCommand := fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDbName)
	db.Exec(dropDatabaseCommand)

	createDatabaseCommand := fmt.Sprintf("CREATE DATABASE %s", testDbName)
	db.Exec(createDatabaseCommand)

	dbConf.DbName = testDbName

	return GetInitPostgreSQLDB(dbConf)
}
