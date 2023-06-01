package main

import (
	"banduslib/internal/api"
	"banduslib/internal/bww"
	"banduslib/internal/database"
	"banduslib/internal/utils"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func initConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("banduslib")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	utils.SetupConsoleLogger()

	config, err := initConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed init configuration")
	}

	var db *gorm.DB
	db, err = database.GetInitSqliteDb(config.SqliteDbPath)
	if err != nil {
		panic(fmt.Sprintf("failed initializing database: %s", err.Error()))
	}

	bwwParser := bww.NewBwwParser()
	bwwFileTuneSplitter := bww.NewBwwFileTuneSplitter()
	dbService := database.NewDbDataService(db)
	apiHandler := api.NewApiHandler(dbService, bwwParser, bwwFileTuneSplitter)
	apiRouter := api.NewApiRouter(apiHandler)

	router := apiRouter.GetEngine()

	log.Info().Msgf("listening on %s", config.ServerUrl)
	log.Fatal().Err(router.Run(config.ServerUrl))
}
