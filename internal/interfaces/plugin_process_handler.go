package interfaces

import "github.com/tomvodi/limepipes-plugin-api/plugin/v1/interfaces"

type PluginProcessHandler interface {
	RunPlugin(pluginID string, executable string) error
	GetPlugin(pluginID string) (interfaces.LimePipesPlugin, error)
	KillPlugins() error
}
