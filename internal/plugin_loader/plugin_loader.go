package plugin_loader

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-plugin"
	wzerolog "github.com/jkratz55/konsul/log/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/common"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/grpc_plugin"
	plugininterfaces "github.com/tomvodi/limepipes-plugin-api/plugin/v1/interfaces"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"github.com/tomvodi/limepipes/internal/interfaces"
	"golang.org/x/exp/maps"
	"os"
	"os/exec"
	"path/filepath"
)

type loader struct {
	pluginClients map[string]*plugin.Client
	plugins       map[string]plugininterfaces.LimePipesPlugin
	pluginInfos   map[string]*messages.PluginInfoResponse
}

func (l *loader) UnloadPlugins() error {
	for pId, client := range l.pluginClients {
		delete(l.plugins, pId)

		client.Kill()
		delete(l.pluginClients, pId)
	}

	return nil
}

func (l *loader) LoadedPlugins() []*messages.PluginInfoResponse {
	return maps.Values(l.pluginInfos)
}

func (l *loader) LoadPluginsFromDir(
	pluginsDir string,
) error {
	allPlugins := l.buildPluginSet()
	for plugId, plug := range allPlugins {
		err := l.loadPlugin(pluginsDir, plugId, plug)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *loader) buildPluginSet() plugin.PluginSet {
	allPlugins := plugin.PluginSet{}
	allPlugins["bww"] = grpc_plugin.NewGrpcPlugin(nil)

	return allPlugins
}

func (l *loader) loadPlugin(
	pluginDir string,
	pluginId string,
	goPlug plugin.Plugin,
) error {
	hcLogger := wzerolog.Wrap(log.Logger)
	pluginExeName := fmt.Sprintf("limepipes-plugin-%s", pluginId)
	pluginExePath := filepath.Join(pluginDir, pluginExeName)

	if _, err := os.Stat("/path/to/whatever"); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("plugin executable '%s' does not exist", pluginExePath)
	}

	clientConf := &plugin.ClientConfig{
		HandshakeConfig: common.HandshakeConfig,
		Plugins: plugin.PluginSet{
			pluginId: goPlug,
		},
		Cmd:              exec.Command("sh", "-c", pluginExePath),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           hcLogger,
	}

	client := plugin.NewClient(clientConf)

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(pluginId)
	if err != nil {
		return err
	}

	lpPlugin, ok := raw.(plugininterfaces.LimePipesPlugin)
	if !ok {
		return fmt.Errorf("plugin %s is not of type LimePipesPlugin", pluginId)
	}

	pInfo, err := lpPlugin.PluginInfo()
	if err != nil {
		return fmt.Errorf("failed getting plugin info from '%s': %v", pluginId, err)
	}

	l.pluginClients[pluginId] = client
	l.plugins[pluginId] = lpPlugin
	l.pluginInfos[pluginId] = pInfo

	return nil
}

func (l *loader) PluginForFileExtension(
	fileExtension string,
) (plugininterfaces.LimePipesPlugin, error) {
	for s, pInfo := range l.pluginInfos {
		for _, ext := range pInfo.FileTypes {
			if ext == fileExtension {
				return l.plugins[s], nil
			}
		}
	}

	return nil, fmt.Errorf("no plugin found for file extension '%s'", fileExtension)
}

func NewPluginLoader() interfaces.PluginLoader {
	return &loader{}
}
