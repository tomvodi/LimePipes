package cmd

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/tomvodi/limepipes/internal/config"
	"github.com/tomvodi/limepipes/internal/interfaces"
	"github.com/tomvodi/limepipes/internal/utils"
)

func NewImportCmd(opts *Options) *cobra.Command {
	importCmd := &cobra.Command{
		Use:   "import [paths...]",
		Short: "Import given files into database",
		Long: `Given files will be parsed and stored in the database. 
When given directory paths, it will import all files of that directory. It will also include 
subdirectories when given the recursive flag.
If a given file that has an extension which is not in the import-file-types, it will be ignored. 
`,
		Args: cobra.MinimumNArgs(1),
		RunE: newImportRunFunc(opts),
	}

	addVerbose(importCmd, opts)
	addDryRun(importCmd, opts)
	addRecursive(importCmd, opts)
	addImportFileTypes(importCmd, opts)
	addSkipFailedFiles(importCmd, opts)

	return importCmd
}

func newImportRunFunc(opts *Options) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, paths []string) error {
		utils.SetupConsoleLogger()
		afs := afero.NewOsFs()

		cfg, err := config.Init()
		if err != nil {
			return fmt.Errorf("failed init configuration: %s", err.Error())
		}

		var pluginLoader interfaces.PluginLoader
		pluginLoader, err = setupPluginLoader(cfg)
		if err != nil {
			return fmt.Errorf("failed setting up plugin loader: %s", err.Error())
		}
		defer pluginLoaderUnload(pluginLoader)

		var dbService interfaces.DataService
		dbService, err = setupDbService(cfg.DbConfig())
		if err != nil {
			return fmt.Errorf("failed setting up database service: %s", err.Error())
		}

		fp := NewFileProcessor(afs, pluginLoader, dbService)

		return fp.ProcessFiles(
			&ProcessFilesOptions{
				ArgPaths:        paths,
				MoveToOutputDir: false,
				ImportToDb:      true,
			},
			opts,
		)
	}
}
