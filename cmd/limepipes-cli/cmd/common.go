package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/tomvodi/limepipes/cmd/limepipes-cli/importtype"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var (
	argRecursive       bool
	argDryRun          bool
	argImportTypes     []string
	argSkipFailedFiles bool
)

func addImportFileTypes(cmd *cobra.Command) {

	cmd.Flags().StringSliceVarP(&argImportTypes, "import-file-types", "i", []string{importtype.BWW.String()},
		fmt.Sprintf("only import files of a specific type. Can be specified with a comma separated list like "+
			"-i=bww,xml. Valid file types are (%s).", strings.Join(importtype.TypeStrings(), ", ")),
	)
}

func addRecursive(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&argRecursive, "recursive", "r", false,
		"import files recursively for all given directory paths",
	)
}

func addDryRun(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&argDryRun, "dry-run", "d", false,
		"import files but only convert them and print infos. Don't import them to the database and don't stop on parsing errors.",
	)
}

func addVerbose(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&argVerbose, "verbose", "v", false,
		"verbose output",
	)
}

func addSkipFailedFiles(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&argSkipFailedFiles, "skip-failed", "s", false,
		"only print error when file couldn't be parsed. If flag is false, the command returns an error on first failed parsing",
	)
}

func invalidImportTypesPassed() []string {
	var invalidImportTypes []string
	for _, aImportType := range argImportTypes {
		if !slices.Contains(importtype.TypeStrings(), aImportType) {
			invalidImportTypes = append(invalidImportTypes, aImportType)
		}
	}

	return invalidImportTypes
}

// getFilesFromDir returns the absolute paths from all files passed in the paths slice.
// If the path is a directory, it will return all files from that directory.
// If the path is a file and has the correct extension, it will return the file.
func getAllFilesFromPaths(paths []string) ([]string, error) {
	var allFiles []string
	for _, p := range paths {
		fp, err := getFilesFromPath(p)
		if err != nil {
			return nil, err
		}
		allFiles = append(allFiles, fp...)
	}

	return allFiles, nil
}

func getFilesFromPath(path string) ([]string, error) {
	var allFiles []string
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// file without valid extension will be ignored
	if !stat.IsDir() && !fileHasValidExtension(stat.Name()) && argVerbose {
		log.Info().Msgf("file %s will be ignored because it doesn't have a valid"+
			"extension (%s)", stat.Name(), strings.Join(argImportTypes, ", "))
		return allFiles, nil
	}

	// file with valid extension will be added to the list
	if !stat.IsDir() && fileHasValidExtension(stat.Name()) {
		allFiles = append(allFiles, path)
		return allFiles, nil
	}

	// if path is a directory, add all files with valid extension to the list
	var dirFiles []string
	if argRecursive {
		dirFiles, err = getFilesRecursively(path)
	} else {
		dirFiles, err = getFilesFromDir(path)
	}
	if err != nil {
		return nil, err
	}
	allFiles = append(allFiles, dirFiles...)

	return allFiles, nil
}

// getFileFromDir returns the absolute paths from all files in the given directory.
func getFilesFromDir(path string) ([]string, error) {
	var filePaths []string
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

	return filePaths, nil
}

// getFilesRecursively returns the absolute paths from all files in the given
// directory and all subdirectories.
func getFilesRecursively(path string) ([]string, error) {
	var filePaths []string

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

	return filePaths, nil
}

func fileHasValidExtension(path string) bool {
	ext := filepath.Ext(path)
	ext = strings.Trim(ext, ".")
	for _, importType := range argImportTypes {
		if ext == importType {
			return true
		}
	}

	return false
}

func checkForInvalidImportTypes() error {
	invalidTypes := invalidImportTypesPassed()
	if len(invalidTypes) > 0 {
		return fmt.Errorf(
			"invalid import file types: %s",
			strings.Join(invalidTypes, ", "),
		)
	}
	return nil
}
