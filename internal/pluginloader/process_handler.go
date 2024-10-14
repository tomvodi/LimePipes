package pluginloader

import (
	"fmt"
	"github.com/hashicorp/go-plugin"
	wzerolog "github.com/jkratz55/konsul/log/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/common"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/grpcplugin"
	plugininterfaces "github.com/tomvodi/limepipes-plugin-api/plugin/v1/interfaces"
	lpcommon "github.com/tomvodi/limepipes/internal/common"
	"os/exec"
)

type ProcessHandler struct {
	pluginSet     plugin.PluginSet
	pluginClients map[string]*plugin.Client
	plugins       map[string]plugininterfaces.LimePipesPlugin
}

func (p *ProcessHandler) KillPlugins() error {
	for pID, client := range p.pluginClients {
		delete(p.plugins, pID)

		client.Kill()
		delete(p.pluginClients, pID)
	}

	return nil
}

// GetPlugin returns a plugin by its ID
// nolint: ireturn
// linter exception is ok here, as the PluginProcessHandler
// interface returns an interface here
func (p *ProcessHandler) GetPlugin(
	pluginID string,
) (plugininterfaces.LimePipesPlugin, error) {
	lp, ok := p.plugins[pluginID]
	if !ok {
		return nil, fmt.Errorf("lp %s not found", pluginID)
	}

	return lp, nil
}

func (p *ProcessHandler) RunPlugin(
	pluginID string,
	executablePath string,
) error {
	hcLogger := wzerolog.Wrap(log.Logger)

	var plug *plugin.Plugin
	ap := p.pluginSet
	for pID, p := range ap {
		if pID == pluginID {
			plug = &p
			break
		}
	}
	if plug == nil {
		return fmt.Errorf("plugin %s not found", pluginID)
	}

	clientConf := &plugin.ClientConfig{
		HandshakeConfig:  common.HandshakeConfig,
		Plugins:          p.pluginSet,
		Cmd:              exec.Command("sh", "-c", executablePath),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           hcLogger,
	}

	client := plugin.NewClient(clientConf)
	p.pluginClients[pluginID] = client

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

	p.plugins[pluginID] = lpPlugin

	return nil
}

func buildPluginSet(
	supportedPlugins lpcommon.PluginList,
) plugin.PluginSet {
	allPlugins := plugin.PluginSet{}
	for _, pID := range supportedPlugins {
		allPlugins[pID] = grpcplugin.NewGrpcPlugin(nil)
	}

	return allPlugins
}

func NewProcessHandler(
	supportedPlugins lpcommon.PluginList,
) *ProcessHandler {
	return &ProcessHandler{
		pluginSet:     buildPluginSet(supportedPlugins),
		pluginClients: make(map[string]*plugin.Client),
		plugins:       make(map[string]plugininterfaces.LimePipesPlugin),
	}
}
