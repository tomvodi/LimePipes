package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
	"github.com/tomvodi/limepipes/internal/api"
	"github.com/tomvodi/limepipes/internal/config"
	"github.com/tomvodi/limepipes/internal/database"
	"github.com/tomvodi/limepipes/internal/interfaces"
	"github.com/tomvodi/limepipes/internal/pluginloader"
)

func setupPluginLoader(
	afs afero.Fs,
	cfg *config.Config,
) (*pluginloader.Loader, error) {
	// TODO: Load plugins from config
	LoadPlugins := []string{
		fileformat.Format_BWW.String(),
	}

	var pluginProcHandler interfaces.PluginProcessHandler = pluginloader.NewProcessHandler(LoadPlugins)
	pluginLoader := pluginloader.NewPluginLoader(
		afs,
		pluginProcHandler,
		LoadPlugins,
	)
	err := pluginLoader.LoadPluginsFromDir(cfg.PluginsDirectoryPath)
	if err != nil {
		return nil, fmt.Errorf("failed loading plugins: %w", err)
	}

	return pluginLoader, nil
}

func pluginLoaderUnload(pluginLoader interfaces.PluginLoader) {
	err := pluginLoader.UnloadPlugins()
	if err != nil {
		log.Fatal().Err(err).Msg("failed unloading plugins")
	}
}

func setupDbService(
	cfg config.DbConfig,
) (*database.Service, error) {
	db, err := database.GetInitPostgreSQLDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed initializing database: %s", err.Error())
	}
	ginValidator := api.NewGinValidator()
	apiModelValidator := api.NewAPIModelValidator(ginValidator)
	return database.NewDbDataService(db, apiModelValidator), nil
}
