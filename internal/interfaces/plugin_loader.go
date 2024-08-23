package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/file_type"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/interfaces"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
)

type PluginLoader interface {
	LoadPluginsFromDir(pluginsDir string) error
	UnloadPlugins() error
	LoadedPlugins() []*messages.PluginInfoResponse
	PluginForFileExtension(fileExtension string) (interfaces.LimePipesPlugin, error)
	FileTypeForFileExtension(fileExtension string) (file_type.Type, error)
}
