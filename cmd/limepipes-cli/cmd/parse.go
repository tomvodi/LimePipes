package cmd

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/tomvodi/limepipes/internal/config"
	"github.com/tomvodi/limepipes/internal/interfaces"
	"github.com/tomvodi/limepipes/internal/utils"
)

func NewParseCmd(opts *Options) *cobra.Command {
	parseCmd := &cobra.Command{
		Use:   "parse [paths...]",
		Short: "Tests parsing given bww files",
		Long: `A bww parser testing command. Parses bww files and moves successful parses into an output directory.

When given directory paths, it will import all files of that directory. It will also include 
subdirectories when given the recursive flag.
If a given file that has an extension which is not in the import-file-types, it will be ignored.`,
		RunE: newParseRunFunc(opts),
	}

	addVerbose(parseCmd, opts)
	addDryRun(parseCmd, opts)
	addRecursive(parseCmd, opts)
	addImportFileTypes(parseCmd, opts)
	addSkipFailedFiles(parseCmd, opts)
	addOutputDir(parseCmd, opts)

	return parseCmd
}

func newParseRunFunc(opts *Options) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, paths []string) error {
		utils.SetupConsoleLogger()
		afs := afero.NewOsFs()

		cfg, err := config.Init()
		if err != nil {
			return fmt.Errorf("failed init configuration: %s", err.Error())
		}

		var pluginLoader interfaces.PluginLoader
		pluginLoader, err = setupPluginLoader(afs, cfg)
		if err != nil {
			return fmt.Errorf("failed setting up plugin loader: %s", err.Error())
		}
		defer pluginLoaderUnload(pluginLoader)

		fp := NewFileProcessor(afs, pluginLoader, nil)

		return fp.ProcessFiles(
			&ProcessFilesOptions{
				ArgPaths:        paths,
				MoveToOutputDir: true,
				ImportToDb:      false,
			},
			opts,
		)
	}
}
