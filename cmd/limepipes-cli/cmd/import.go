package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/tomvodi/limepipes/internal/api"
	"github.com/tomvodi/limepipes/internal/common"
	"github.com/tomvodi/limepipes/internal/config"
	"github.com/tomvodi/limepipes/internal/database"
	"github.com/tomvodi/limepipes/internal/interfaces"
	"github.com/tomvodi/limepipes/internal/pluginloader"
	"github.com/tomvodi/limepipes/internal/utils"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

var (
	argVerbose bool
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import [paths...]",
	Short: "Import given files into database",
	Long: `Given files will be parsed and stored in the database. 
When given directory paths, it will import all files of that directory. It will also include 
subdirectories when given the recursive flag.
If a given file that has an extension which is not in the import-file-types, it will be ignored. 
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, paths []string) error {
		utils.SetupConsoleLogger()
		fs := afero.NewOsFs()

		cfg, err := config.Init()
		if err != nil {
			return fmt.Errorf("failed init configuration: %s", err.Error())
		}

		// TODO: Load plugins from config
		LoadPlugins := []string{
			"bww",
		}

		var pluginProcHandler interfaces.PluginProcessHandler = pluginloader.NewProcessHandler(LoadPlugins)
		var pluginLoader interfaces.PluginLoader = pluginloader.NewPluginLoader(
			fs,
			pluginProcHandler,
			LoadPlugins,
		)
		err = pluginLoader.LoadPluginsFromDir(cfg.PluginsDirectoryPath)
		if err != nil {
			log.Fatal().Err(err).Msg("failed loading plugins")
		}
		defer func(pluginLoader interfaces.PluginLoader) {
			err := pluginLoader.UnloadPlugins()
			if err != nil {
				log.Fatal().Err(err).Msg("failed unloading plugins")
			}
		}(pluginLoader)

		var db *gorm.DB
		db, err = database.GetInitPostgreSQLDB(cfg.DbConfig())
		if err != nil {
			return fmt.Errorf("failed initializing database: %s", err.Error())
		}
		ginValidator := api.NewGinValidator()
		apiModelValidator := api.NewAPIModelValidator(ginValidator)
		dbService := database.NewDbDataService(db, apiModelValidator)

		err = checkForInvalidImportTypes()
		if err != nil {
			return err
		}

		allFiles, err := getAllFilesFromPaths(paths)
		if err != nil {
			return fmt.Errorf("failed getting files: %s", err.Error())
		}

		log.Info().Msgf("found %d files for import", len(allFiles))

		allFileCnt := len(allFiles)
		for i, file := range allFiles {
			fExt := filepath.Ext(file)
			if fExt == "" {
				if argSkipFailedFiles {
					log.Error().Err(err).Msgf("import file %s does not have an extension", file)
					continue
				}

				return err
			}

			filePlugin, err := pluginLoader.PluginForFileExtension(".bww")
			if err != nil {
				if argSkipFailedFiles {
					log.Error().Err(err).Msgf("failed getting plugin for file %s with extension %s",
						file, fExt)
					continue
				}

				return err
			}

			fileData, err := os.ReadFile(file)
			if err != nil {
				if argSkipFailedFiles {
					log.Error().Err(err).Msgf("failed reading file %s", file)
					continue
				}

				return err
			}

			log.Info().Msgf("importing file %d/%d %s", i+1, allFileCnt, file)

			tunesImport, err := filePlugin.Import(fileData)
			if err != nil {
				if argSkipFailedFiles {
					log.Error().Err(err).Msgf("failed parsing file %s", file)
					continue
				}

				return fmt.Errorf("failed parsing file %s: %v", file, err)
			}

			if argVerbose {
				log.Info().Msgf("(%d/%d) successfully parsed %d tunes from file %s",
					i+1,
					allFileCnt,
					len(tunesImport.ImportedTunes),
					file,
				)
			}

			fType, err := pluginLoader.FileTypeForFileExtension(fExt)
			if err != nil {
				if argSkipFailedFiles {
					log.Error().Err(err).Msgf("failed getting file type for file %s with extension %s",
						file, fExt)
					continue
				}

				return err
			}

			fInfo, err := common.NewImportFileInfoFromLocalFile(file, fType)
			if err != nil {
				if argSkipFailedFiles {
					log.Error().Err(err).Msg("failed creating import file info")
					continue
				}

				return err
			}

			_, _, err = dbService.ImportTunes(tunesImport.ImportedTunes, fInfo)
			if err != nil {
				if argSkipFailedFiles {
					log.Error().Err(err).Msg("failed importing tunes")
					continue
				}

				return fmt.Errorf("failed importing tunes: %s", err.Error())
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	addVerbose(importCmd)
	addDryRun(importCmd)
	addRecursive(importCmd)
	addImportFileTypes(importCmd)
	addSkipFailedFiles(importCmd)
}
