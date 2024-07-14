package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes/internal/api"
	"github.com/tomvodi/limepipes/internal/bww"
	"github.com/tomvodi/limepipes/internal/common/music_model/helper"
	"github.com/tomvodi/limepipes/internal/config"
	"github.com/tomvodi/limepipes/internal/database"
	"github.com/tomvodi/limepipes/internal/utils"
	"gorm.io/gorm"
)

func main() {
	utils.SetupConsoleLogger()

	cfg, err := config.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("failed init configuration")
	}

	var db *gorm.DB
	db, err = database.GetInitSqliteDb(cfg.SqliteDbPath)
	if err != nil {
		panic(fmt.Sprintf("failed initializing database: %s", err.Error()))
	}

	bwwParser := bww.NewBwwParser()
	bwwFileTuneSplitter := bww.NewBwwFileTuneSplitter()
	dbService := database.NewDbDataService(db)
	tuneFixer := helper.NewTuneFixer()
	apiHandler := api.NewApiHandler(dbService, bwwParser, bwwFileTuneSplitter, tuneFixer)
	apiRouter := api.NewApiRouter(apiHandler)

	router := apiRouter.GetEngine()

	log.Info().Msgf("listening on %s", cfg.ServerUrl)
	log.Fatal().Err(router.RunTLS(cfg.ServerUrl, cfg.TlsCertPath, cfg.TlsCertKeyPath))
}
