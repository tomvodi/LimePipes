//go:build wireinject
// +build wireinject

package initialize

import (
	"github.com/google/wire"
	"github.com/spf13/afero"
	"github.com/tomvodi/limepipes/internal/api"
	"github.com/tomvodi/limepipes/internal/common"
	"github.com/tomvodi/limepipes/internal/config"
	"github.com/tomvodi/limepipes/internal/database"
	"github.com/tomvodi/limepipes/internal/health"
	"github.com/tomvodi/limepipes/internal/interfaces"
	"github.com/tomvodi/limepipes/internal/pluginloader"
	"gorm.io/gorm"
)

func PluginLoader(
	pluginList common.PluginList,
) *pluginloader.Loader {
	wire.Build(
		wire.InterfaceValue(new(afero.Fs), afero.NewOsFs()),
		wire.Bind(new(interfaces.PluginProcessHandler), new(*pluginloader.ProcessHandler)),
		pluginloader.NewProcessHandler,
		pluginloader.NewPluginLoader,
	)

	return &pluginloader.Loader{}
}

func ApiHandler(
	db *gorm.DB,
	healthConfig config.HealthConfig,
	pluginloader interfaces.PluginLoader,
) (*api.Handler, error) {
	wire.Build(
		api.NewGinValidator,
		api.NewAPIModelValidator,
		wire.Bind(new(interfaces.APIModelValidator), new(*api.ModelValidator)),
		database.NewDbDataService,
		wire.Bind(new(interfaces.DataService), new(*database.Service)),
		health.NewHealthCheck,
		wire.Bind(new(interfaces.HealthChecker), new(*health.Check)),
		api.NewAPIHandler,
	)

	return &api.Handler{}, nil
}
