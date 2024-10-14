package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
	"github.com/tomvodi/limepipes/internal/apigen"
	"github.com/tomvodi/limepipes/internal/common"
	"github.com/tomvodi/limepipes/internal/config"
	"github.com/tomvodi/limepipes/internal/database"
	"github.com/tomvodi/limepipes/internal/initialize"
	"github.com/tomvodi/limepipes/internal/interfaces"
	"github.com/tomvodi/limepipes/internal/utils"
	"gorm.io/gorm"
	"strings"
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

	// TODO: Load plugins from config
	LoadPlugins := common.PluginList{
		strings.ToLower(fileformat.Format_BWW.String()),
	}
	var pluginLoader interfaces.PluginLoader = initialize.PluginLoader(LoadPlugins)

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

	apiHandler, err := initialize.ApiHandler(db, cfg.HealthConfig(), pluginLoader)
	if err != nil {
		panic(fmt.Sprintf("failed initializing health check: %s", err.Error()))
	}

	engine := setupGinEngine()
	router := apigen.NewRouterWithGinEngine(
		engine,
		apigen.ApiHandleFunctions{
			ApiHandler: apiHandler,
		},
	)

	log.Info().Msgf("listening on %s", cfg.ServerURL)
	log.Fatal().Err(router.RunTLS(cfg.ServerURL, cfg.TLSCertPath, cfg.TLSCertKeyPath))
}
