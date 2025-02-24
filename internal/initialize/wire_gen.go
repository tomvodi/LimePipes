// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package initialize

import (
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

// Injectors from wire.go:

func PluginLoader(pluginList common.PluginList) *pluginloader.Loader {
	fs := _wireFsValue
	processHandler := pluginloader.NewProcessHandler(pluginList)
	loader := pluginloader.NewPluginLoader(fs, processHandler, pluginList)
	return loader
}

var (
	_wireFsValue = afero.NewOsFs()
)

func ApiHandler(db *gorm.DB, healthConfig config.HealthConfig, pluginloader2 interfaces.PluginLoader) (*api.Handler, error) {
	validate := api.NewGinValidator()
	modelValidator := api.NewAPIModelValidator(validate)
	service := database.NewDbDataService(db, modelValidator)
	check, err := health.NewHealthCheck(healthConfig, db)
	if err != nil {
		return nil, err
	}
	handler := api.NewAPIHandler(service, pluginloader2, check)
	return handler, nil
}
