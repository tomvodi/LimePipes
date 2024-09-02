package pluginloader

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-plugin"
	wzerolog "github.com/jkratz55/konsul/log/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/common"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/file_type"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/grpc_plugin"
	plugininterfaces "github.com/tomvodi/limepipes-plugin-api/plugin/v1/interfaces"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"golang.org/x/exp/maps"
	"os"
	"os/exec"
	"path/filepath"
)

type Loader struct {
	pluginClients map[string]*plugin.Client
	plugins       map[string]plugininterfaces.LimePipesPlugin
	pluginInfos   map[string]*messages.PluginInfoResponse
}

func (l *Loader) FileTypeForFileExtension(fileExtension string) (file_type.Type, error) {
	if len(l.pluginInfos) == 0 {
		return file_type.Type_Unknown, errors.New("no plugins loaded")
	}

	for _, pInfo := range l.pluginInfos {
		for _, ext := range pInfo.FileExtensions {
			if ext == fileExtension {
				return pInfo.FileType, nil
			}
		}
	}

	return file_type.Type_Unknown,
		fmt.Errorf("no plugin found that handles file extension '%s'", fileExtension)
}

func (l *Loader) UnloadPlugins() error {
	for pID, client := range l.pluginClients {
		delete(l.plugins, pID)

		client.Kill()
		delete(l.pluginClients, pID)
	}

	return nil
}

func (l *Loader) LoadedPlugins() []*messages.PluginInfoResponse {
	return maps.Values(l.pluginInfos)
}

func (l *Loader) LoadPluginsFromDir(
	pluginsDir string,
) error {
	allPlugins := l.buildPluginSet()
	for pID, plug := range allPlugins {
		err := l.loadPlugin(pluginsDir, pID, plug)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Loader) buildPluginSet() plugin.PluginSet {
	allPlugins := plugin.PluginSet{}
	allPlugins["bww"] = grpc_plugin.NewGrpcPlugin(nil)

	return allPlugins
}

func (l *Loader) loadPlugin(
	pluginDir string,
	pluginID string,
	goPlug plugin.Plugin,
) error {
	hcLogger := wzerolog.Wrap(log.Logger)
	pluginExeName := fmt.Sprintf("limepipes-plugin-%s", pluginID)
	pluginExePath := filepath.Join(pluginDir, pluginExeName)

	if _, err := os.Stat(pluginExePath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("plugin executable '%s' does not exist", pluginExePath)
	}

	clientConf := &plugin.ClientConfig{
		HandshakeConfig: common.HandshakeConfig,
		Plugins: plugin.PluginSet{
			pluginID: goPlug,
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
	raw, err := rpcClient.Dispense(pluginID)
	if err != nil {
		return err
	}

	lpPlugin, ok := raw.(plugininterfaces.LimePipesPlugin)
	if !ok {
		return fmt.Errorf("plugin %s is not of type LimePipesPlugin", pluginID)
	}

	pInfo, err := lpPlugin.PluginInfo()
	if err != nil {
		return fmt.Errorf("failed getting plugin info from '%s': %v", pluginID, err)
	}

	l.pluginClients[pluginID] = client
	l.plugins[pluginID] = lpPlugin
	l.pluginInfos[pluginID] = pInfo

	return nil
}

// PluginForFileExtension returns the plugin that can handle the given file extension.
// nolint: ireturn
// linter exception is ok here, as the PluginLoader interface returns an interface here
func (l *Loader) PluginForFileExtension(
	fileExtension string,
) (plugininterfaces.LimePipesPlugin, error) {
	for s, pInfo := range l.pluginInfos {
		for _, ext := range pInfo.FileExtensions {
			if ext == fileExtension {
				return l.plugins[s], nil
			}
		}
	}

	return nil, fmt.Errorf("no plugin found for file extension '%s'", fileExtension)
}

func NewPluginLoader() *Loader {
	return &Loader{
		pluginClients: map[string]*plugin.Client{},
		plugins:       map[string]plugininterfaces.LimePipesPlugin{},
		pluginInfos:   map[string]*messages.PluginInfoResponse{},
	}
}
