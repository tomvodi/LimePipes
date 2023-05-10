/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"banduslib/internal/bww"
	"banduslib/internal/common/music_model"
	"banduslib/internal/utils"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	ImportTypeBww string = "bww"
)

var allImportTypes = []string{
	ImportTypeBww,
}

var (
	verbose     bool
	recursive   bool
	dryRun      bool
	importTypes []string
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import [paths...]",
	Short: "Import given files into database",
	Long: `Given files will internally be converted to MusicXML and stored in the database. 
When given directory paths, it will import all files of that directory. It will also include 
subdirectories when given the recursive flag.
If a given file that has an extension which is not in the import-file-types, it will be ignored. 
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		utils.SetupConsoleLogger()
		invalidTypes := hasInvalidImportTypes()
		if len(invalidTypes) > 0 {
			return fmt.Errorf(
				"invalid import file types: %s",
				strings.Join(invalidTypes, ", "),
			)
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

		parser := bww.NewBwwParser()
		for _, file := range allFiles {
			fileData, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			if verbose {
				log.Info().Msgf("importing file %s", file)
			}
			var bwwTunes []*music_model.Tune
			bwwTunes, err = parser.ParseBwwData(fileData)
			if err != nil {
				return fmt.Errorf("faied parsing file %s: %v", file, err)
			}
			if verbose {
				log.Info().Msgf("successfully parsed %d tunes from file %s",
					len(bwwTunes),
					file,
				)
			}
		}

		fmt.Printf("import called with %v, importing files of type %s", args, strings.Join(importTypes, ", "))
		return nil
	},
}

func hasInvalidImportTypes() []string {
	var invalidImportTypes []string
	for _, importType := range importTypes {
		isValidImportType := false
		for _, allImportType := range allImportTypes {
			if allImportType == importType {
				isValidImportType = true
				break
			}
		}
		if !isValidImportType {
			invalidImportTypes = append(invalidImportTypes, importType)
		}
	}

	return invalidImportTypes
}

func getAllFilesFromArgs(args []string) ([]string, error) {
	var allFiles []string
	for _, arg := range args {
		stat, err := os.Stat(arg)
		if err != nil {
			return nil, err
		}
		if stat.IsDir() {
			dirFiles, err := getFilesFromDir(arg)
			if err != nil {
				return nil, err
			}
			allFiles = append(allFiles, dirFiles...)
		} else {
			if fileHasValidExtension(stat.Name()) {
				allFiles = append(allFiles, arg)
			} else {
				if verbose {
					log.Info().Msgf("file %s will be ignored because it doesn't have a valid"+
						"extension (%s)", strings.Join(importTypes, ", "))
				}
			}
		}
	}

	return allFiles, nil
}

func getFilesFromDir(path string) ([]string, error) {
	var filePaths []string
	if recursive {
		filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if fileHasValidExtension(path) {
				filePaths = append(filePaths, path)
			}

			return nil
		})
	} else {
		files, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if fileHasValidExtension(file.Name()) {
				filePaths = append(filePaths, filepath.Join(path, file.Name()))
			}
		}
	}

	return filePaths, nil
}

func fileHasValidExtension(path string) bool {
	ext := filepath.Ext(path)
	ext = strings.Trim(ext, ".")
	for _, importType := range importTypes {
		if ext == importType {
			return true
		}
	}

	return false
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"verbose output",
	)
	importCmd.Flags().BoolVarP(&recursive, "recursive", "r", false,
		"import files recursively for all given directory paths",
	)
	importCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false,
		"import files but only convert them and print infos. Don't import them to the database",
	)
	importCmd.Flags().StringSliceVarP(&importTypes, "import-file-types", "i", []string{"bww"},
		fmt.Sprintf("only import files of a specific type. Can be specified with a comma separated list like "+
			"-i=bww,xml. Valid file types are (%s).", strings.Join(allImportTypes, ", ")),
	)
}
