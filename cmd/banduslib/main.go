package main

import (
	"banduslib/internal/api"
	"banduslib/internal/bww"
	"banduslib/internal/common/music_model/helper"
	"banduslib/internal/config"
	"banduslib/internal/database"
	"banduslib/internal/utils"
	"fmt"
	"github.com/rs/zerolog/log"
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
	log.Fatal().Err(router.Run(cfg.ServerUrl))
}
