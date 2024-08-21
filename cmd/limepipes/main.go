package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes/internal/api"
	"github.com/tomvodi/limepipes/internal/api_gen"
	"github.com/tomvodi/limepipes/internal/config"
	"github.com/tomvodi/limepipes/internal/database"
	"github.com/tomvodi/limepipes/internal/health"
	"github.com/tomvodi/limepipes/internal/interfaces"
	"github.com/tomvodi/limepipes/internal/plugin_loader"
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

	pluginLoader := plugin_loader.NewPluginLoader()
	err = pluginLoader.LoadPluginsFromDir(cfg.PluginsDirectoryPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed loading plugins")
	}
	defer func(pluginLoader interfaces.PluginLoader) {
		err := pluginLoader.UnloadPlugins()
		if err != nil {
			log.Fatal().Err(err).Msg("failed unloading plugins")
		}
	}(pluginLoader)

	var db *gorm.DB
	db, err = database.GetInitPostgreSQLDB(cfg.DbConfig())
	if err != nil {
		panic(fmt.Sprintf("failed initializing database: %s", err.Error()))
	}

	ginValidator := api.NewGinValidator()
	apiModelValidator := api.NewApiModelValidator(ginValidator)
	dbService := database.NewDbDataService(db, apiModelValidator)
	healthChecker, err := health.NewHealthCheck(cfg.HealthConfig(), db)
	if err != nil {
		panic(fmt.Sprintf("failed initializing health check: %s", err.Error()))
	}
	apiHandler := api.NewApiHandler(
		dbService,
		pluginLoader,
		healthChecker,
	)
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
