package pluginloader

import (
	"errors"
	"fmt"
	"github.com/spf13/afero"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
	plugininterfaces "github.com/tomvodi/limepipes-plugin-api/plugin/v1/interfaces"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"github.com/tomvodi/limepipes/internal/interfaces"
	"golang.org/x/exp/maps"
	"os"
	"path/filepath"
	"slices"
)

type Loader struct {
	afs              afero.Fs
	processHandler   interfaces.PluginProcessHandler
	supportedPlugins []string
	pluginInfos      map[string]*messages.PluginInfoResponse
}

func (l *Loader) FileFormatForFileExtension(fileExtension string) (fileformat.Format, error) {
	if len(l.pluginInfos) == 0 {
		return fileformat.Format_Unknown, errors.New("no plugins loaded")
	}

	for _, pInfo := range l.pluginInfos {
		for _, ext := range pInfo.FileExtensions {
			if ext == fileExtension {
				return pInfo.FileFormat, nil
			}
		}
	}

	return fileformat.Format_Unknown,
		fmt.Errorf("no plugin found that handles file extension '%s'", fileExtension)
}

func (l *Loader) UnloadPlugins() error {
	return l.processHandler.KillPlugins()
}

func (l *Loader) LoadedPlugins() []*messages.PluginInfoResponse {
	return maps.Values(l.pluginInfos)
}

func (l *Loader) LoadPluginsFromDir(
	pluginsDir string,
) error {
	for _, pID := range l.supportedPlugins {
		err := l.loadPlugin(pluginsDir, pID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Loader) loadPlugin(
	pluginDir string,
	pluginID string,
) error {
	pluginExeName := fmt.Sprintf("limepipes-plugin-%s", pluginID)
	pluginExePath := filepath.Join(pluginDir, pluginExeName)

	if _, err := l.afs.Stat(pluginExePath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("plugin executable '%s' does not exist", pluginExePath)
	}

	err := l.processHandler.RunPlugin(pluginID, pluginExePath)
	if err != nil {
		return fmt.Errorf("failed running plugin %s: %v", pluginID, err)
	}

	lpPlugin, err := l.processHandler.GetPlugin(pluginID)
	if err != nil {
		return fmt.Errorf("failed getting plugin %s: %v", pluginID, err)
	}

	pInfo, err := lpPlugin.PluginInfo()
	if err != nil {
		return fmt.Errorf("failed getting plugin info from '%s': %v", pluginID, err)
	}

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
		if slices.Contains(pInfo.FileExtensions, fileExtension) {
			lp, err := l.processHandler.GetPlugin(s)
			if err != nil {
				return nil, err
			}

			return lp, nil
		}
	}

	return nil, fmt.Errorf("no plugin found for file extension '%s'", fileExtension)
}

func (l *Loader) FileExtensionsForFileFormat(
	format fileformat.Format,
) ([]string, error) {
	for _, pInfo := range l.pluginInfos {
		if pInfo.FileFormat == format {
			return pInfo.FileExtensions, nil
		}
	}

	return nil, fmt.Errorf("no plugin found for file format '%s'", format.String())
}

func NewPluginLoader(
	afs afero.Fs,
	processHandler interfaces.PluginProcessHandler,
	supportedPlugins []string,
) *Loader {
	return &Loader{
		afs:              afs,
		processHandler:   processHandler,
		supportedPlugins: supportedPlugins,
		pluginInfos:      map[string]*messages.PluginInfoResponse{},
	}
}
