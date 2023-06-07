package cmd

import (
	"banduslib/internal/bww"
	"banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/expander"
	"banduslib/internal/common/music_model/helper"
	"banduslib/internal/config"
	"banduslib/internal/database"
	"banduslib/internal/utils"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"os"
)

var (
	verbose bool
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
	RunE: func(cmd *cobra.Command, args []string) error {
		utils.SetupConsoleLogger()
		cfg, err := config.Init()
		if err != nil {
			return fmt.Errorf("failed init configuration: %s", err.Error())
		}

		var db *gorm.DB
		db, err = database.GetInitSqliteDb(cfg.SqliteDbPath)
		if err != nil {
			return fmt.Errorf("failed initializing database: %s", err.Error())
		}
		dbService := database.NewDbDataService(db)
		bwwFileTuneSplitter := bww.NewBwwFileTuneSplitter()
		tuneFixer := helper.NewTuneFixer()

		err = checkForInvalidImportTypes()
		if err != nil {
			return err
		}

		allFiles, err := getAllFilesFromArgs(args)
		if err != nil {
			return fmt.Errorf("failed getting files: %s", err.Error())
		}

		log.Info().Msgf("found %d files for import", len(allFiles))

		importedTunesCnt := 0
		embExpander := expander.NewEmbellishmentExpander()
		parser := bww.NewBwwParser(embExpander)
		allFileCnt := len(allFiles)
		for i, file := range allFiles {
			fileData, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			log.Info().Msgf("importing file %d/%d %s", i+1, allFileCnt, file)
			var muModel music_model.MusicModel
			muModel, err = parser.ParseBwwData(fileData)
			if err != nil {
				if skipFailedFiles {
					log.Error().Err(err).Msgf("failed parsing file %s", file)
					continue
				} else {
					return fmt.Errorf("failed parsing file %s: %v", file, err)
				}
			}

			if verbose {
				log.Info().Msgf("(%d/%d) successfully parsed %d tunes from file %s",
					i+1,
					allFileCnt,
					len(muModel),
					file,
				)
			}

			tuneFixer.Fix(muModel)

			bwwFileTuneData, err := bwwFileTuneSplitter.SplitFileData(fileData)
			if err != nil {
				msg := fmt.Sprintf("failed splitting file %s by tunes: %s", file, err.Error())
				if skipFailedFiles {
					log.Error().Msg(msg)
					continue
				} else {
					return fmt.Errorf(msg)
				}
			}

			if len(bwwFileTuneData.TuneTitles()) != len(muModel) {
				log.Warn().Msgf("split bww file and music model don't have the same amount of tunes: %s",
					file)
			}

			fInfo, err := common.NewImportFileInfoFromLocalFile(file)
			if err != nil {
				if skipFailedFiles {
					log.Error().Err(err).Msg("failed creating import file info")
					continue
				} else {
					return err
				}
			}
			apiTunes, err := dbService.ImportMusicModel(muModel, fInfo, bwwFileTuneData)
			if err != nil {
				if skipFailedFiles {
					log.Error().Err(err).Msg("failed importing tunes")
					continue
				} else {
					return fmt.Errorf("failed importing tunes: %s", err.Error())
				}
			}
			for _, tune := range apiTunes {
				if tune.ImportedToDatabase {
					importedTunesCnt++
				}
			}
		}

		log.Info().Msgf("from %d files, imported %d new tunes into database", len(allFiles), importedTunesCnt)

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
