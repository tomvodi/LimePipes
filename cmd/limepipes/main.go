package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes/internal/api"
	"github.com/tomvodi/limepipes/internal/api_gen"
	"github.com/tomvodi/limepipes/internal/bww"
	"github.com/tomvodi/limepipes/internal/common/music_model/helper"
	"github.com/tomvodi/limepipes/internal/config"
	"github.com/tomvodi/limepipes/internal/database"
	"github.com/tomvodi/limepipes/internal/utils"
	"gorm.io/gorm"
)

func setupGinEngine() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://localhost:3000"},
		AllowMethods: []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders: []string{"Origin", "Content-type"},
	}))

	return router
}

func main() {
	utils.SetupConsoleLogger()

	cfg, err := config.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("failed init configuration")
	}

	var db *gorm.DB
	db, err = database.GetInitPostgreSQLDB(cfg.DbConfig())
	if err != nil {
		panic(fmt.Sprintf("failed initializing database: %s", err.Error()))
	}

	bwwParser := bww.NewBwwParser()
	bwwFileTuneSplitter := bww.NewBwwFileTuneSplitter()
	ginValidator := api.NewGinValidator()
	apiModelValidator := api.NewApiModelValidator(ginValidator)
	dbService := database.NewDbDataService(db, apiModelValidator)
	tuneFixer := helper.NewTuneFixer()
	apiHandler := api.NewApiHandler(dbService, bwwParser, bwwFileTuneSplitter, tuneFixer)
	engine := setupGinEngine()
	router := api_gen.NewRouterWithGinEngine(
		engine,
		api_gen.ApiHandleFunctions{
			ApiHandler: apiHandler,
		},
	)

	log.Info().Msgf("listening on %s", cfg.ServerUrl)
	log.Fatal().Err(router.RunTLS(cfg.ServerUrl, cfg.TlsCertPath, cfg.TlsCertKeyPath))
}
