package cmd

import (
	"banduslib/internal/bww"
	"banduslib/internal/common/music_model"
	"banduslib/internal/utils"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
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
		err := checkForInvalidImportTypes()
		if err != nil {
			return err
		}

		allFiles, err := getAllFilesFromArgs(args)
		if err != nil {
			return fmt.Errorf("failed getting files: %s", err.Error())
		}
		if verbose {
			log.Info().Msg("Processing files: ")
			for _, file := range allFiles {
				log.Info().Msg(file)
			}
		}

		tuneMap := make(map[string]*music_model.Tune)
		parser := bww.NewBwwParser()
		allFileCnt := len(allFiles)
		for i, file := range allFiles {
			fileData, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			if verbose {
				log.Info().Msgf("importing file %s", file)
			}
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

			for _, tune := range muModel {
				existingTune, ok := tuneMap[tune.Title]
				if ok && reflect.DeepEqual(tune, existingTune) {
					if verbose {
						log.Info().Msgf("tune with title %s was already parsed and is equals existing one", tune.Title)
					}
				}
				if !ok {
					tuneMap[tune.Title] = tune
				}
			}
		}

		log.Info().Msgf("from %d files, got %d tunes without duplicates", len(allFiles), len(tuneMap))

		fmt.Printf("import called with %v, importing files of type %s", args, strings.Join(importTypes, ", "))
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
