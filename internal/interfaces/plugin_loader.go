package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/interfaces"
)

type PluginLoader interface {
	LoadPluginsFromDir(pluginsDir string) error
	UnloadPlugins() error
	PluginForFileExtension(fileExtension string) (interfaces.LimePipesPlugin, error)
	FileTypeForFileExtension(fileExtension string) (fileformat.Format, error)
}
