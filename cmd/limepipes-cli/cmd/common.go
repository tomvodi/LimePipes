package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	ImportTypeBww string = "bww"
)

var allImportTypes = []string{
	ImportTypeBww,
}

var (
	recursive       bool
	dryRun          bool
	importTypes     []string
	skipFailedFiles bool
)

func addImportFileTypes(cmd *cobra.Command) {

	cmd.Flags().StringSliceVarP(&importTypes, "import-file-types", "i", []string{ImportTypeBww},
		fmt.Sprintf("only import files of a specific type. Can be specified with a comma separated list like "+
			"-i=bww,xml. Valid file types are (%s).", strings.Join(allImportTypes, ", ")),
	)
}

func addRecursive(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false,
		"import files recursively for all given directory paths",
	)
}

func addDryRun(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false,
		"import files but only convert them and print infos. Don't import them to the database and don't stop on parsing errors.",
	)
}

func addVerbose(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"verbose output",
	)
}

func addSkipFailedFiles(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&skipFailedFiles, "skip-failed", "s", false,
		"only print error when file couldn't be parsed. If flag is false, the command returns an error on first failed parsing",
	)
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
						"extension (%s)", stat.Name(), strings.Join(importTypes, ", "))
				}
			}
		}
	}

	return allFiles, nil
}

func getFilesFromDir(path string) ([]string, error) {
	var filePaths []string
	if recursive {
		err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
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
		if err != nil {
			return nil, err
		}
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

func checkForInvalidImportTypes() error {
	invalidTypes := hasInvalidImportTypes()
	if len(invalidTypes) > 0 {
		return fmt.Errorf(
			"invalid import file types: %s",
			strings.Join(invalidTypes, ", "),
		)
	}
	return nil
}
