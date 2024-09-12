package cmd

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"github.com/tomvodi/limepipes/cmd/limepipes-cli/importtype"
	"github.com/tomvodi/limepipes/internal/common"
	"github.com/tomvodi/limepipes/internal/interfaces"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// GetFilesOptions contains options for getting files from directories.
type GetFilesOptions struct {
	Verbose        bool
	Recursive      bool
	FileExtensions []string // only search for files with these extensions
}

// ProcessFilesOptions contains the main options for processing files.
type ProcessFilesOptions struct {
	ArgPaths        []string // List of files and directories passed as arguments
	MoveToOutputDir bool     // Move successful parsed files to output directory
	ImportToDb      bool     // Import successful parsed files to database
}

// ProcessFileContext contains information about the current file being processed.
type ProcessFileContext struct {
	Filepath      string
	CurrentFileNr int // Current file number
	AllFilesCnt   int // Total number of files to process
}

type FileProcessor struct {
	afs afero.Fs
	pl  interfaces.PluginLoader
	ds  interfaces.DataService
}

// revive:disable:cognitive-complexity this function has complexity of 8 which is fine with threshold of 7
func (fp *FileProcessor) ProcessFiles(
	pfo *ProcessFilesOptions,
	opts *Options,
) error {
	err := fp.checkPreconditions(pfo, opts)
	if err != nil {
		return err
	}

	gfo, err := fp.setupGetFilesOptions(opts)
	if err != nil {
		return err
	}

	allFiles, err := getAllFilesFromPaths(fp.afs, pfo.ArgPaths, gfo)
	if err != nil {
		return fmt.Errorf("failed getting files: %s", err.Error())
	}

	pfc := &ProcessFileContext{
		AllFilesCnt: len(allFiles),
	}

	for i, file := range allFiles {
		if opts.Verbose {
			log.Info().Msgf("processing file %d/%d %s",
				pfc.CurrentFileNr,
				pfc.AllFilesCnt,
				pfc.Filepath,
			)
		}

		pfc.Filepath = file
		pfc.CurrentFileNr = i + 1

		err = fp.processFile(pfc, opts, pfo)
		if err != nil {
			return err
		}
	}

	return nil
}

// revive:enable:cognitive-complexity

// checkPreconditionsAndInitConfig checks for:
// - invalid import types
// - creates the output directory if necessary
func (fp *FileProcessor) checkPreconditions(
	pfo *ProcessFilesOptions,
	opts *Options,
) error {
	err := checkForInvalidImportTypes(opts.ImportTypes)
	if err != nil {
		return err
	}

	err = fp.setupOutputDir(pfo, opts)
	if err != nil {
		return err
	}
	return nil
}

// checkForInvalidImportTypes checks if the given import types are valid.
func checkForInvalidImportTypes(iTypes []string) error {
	if len(iTypes) == 0 {
		return fmt.Errorf("no import types passed")
	}

	invalidTypes := invalidImportTypesPassed(iTypes)
	if len(invalidTypes) > 0 {
		return fmt.Errorf(
			"invalid import file types: %s",
			strings.Join(invalidTypes, ", "),
		)
	}
	return nil
}

// setupOutputDir creates the output directory if the successfully parsed files
// should be moved to it. It modifies the output directory to an absolute path if necessary
// and sets the output directory in the options.
func (fp *FileProcessor) setupOutputDir(
	po *ProcessFilesOptions,
	opts *Options,
) error {
	if !po.MoveToOutputDir {
		return nil
	}

	if !filepath.IsAbs(opts.OutputDir) {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		opts.OutputDir = filepath.Join(wd, opts.OutputDir)
	}
	err := fp.afs.MkdirAll(opts.OutputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed creating output directory: %s", err.Error())
	}

	if opts.Verbose {
		log.Info().Msgf("successful parsed files will be moved to: %s", opts.OutputDir)
	}

	return nil
}

// invalidImportTypesPassed takes a list of passed import types and returns
// a list of the import types that are not valid.
// It returns nil if all given import types are valid.
func invalidImportTypesPassed(importTypes []string) []string {
	var invalidImportTypes []string
	validTypes := importtype.AllTypes()
	for _, aImportType := range importTypes {
		if !slices.Contains(validTypes, aImportType) {
			invalidImportTypes = append(invalidImportTypes, aImportType)
		}
	}

	return invalidImportTypes
}

// setupGetFilesOptions returns a GetFilesOptions struct with the given options and
// the valid file extensions for the given import types.
func (fp *FileProcessor) setupGetFilesOptions(
	opts *Options,
) (*GetFilesOptions, error) {
	fileExtensions, err := fp.validFileExtensionsForImportTypes(opts.ImportTypes)
	if err != nil {
		return nil, err
	}

	return &GetFilesOptions{
		Verbose:        opts.Verbose,
		Recursive:      opts.Recursive,
		FileExtensions: fileExtensions,
	}, nil
}

// validFileExtensionsForImportTypes returns all valid file extensions for the given import types that
// are supported by the plugins.
func (fp *FileProcessor) validFileExtensionsForImportTypes(
	iTypes []string,
) ([]string, error) {
	var fileExtensions []string
	mapping := importtype.FileFormatMapping()
	for _, importType := range iTypes {
		ff, ok := mapping[importType]
		if !ok {
			return nil, fmt.Errorf("invalid import type %s", importType)
		}
		typeExt, err := fp.pl.FileExtensionsForFileFormat(ff)
		if err != nil {
			return nil, fmt.Errorf("failed getting file extensions for file format %s: %s", ff.String(), err.Error())
		}

		fileExtensions = append(fileExtensions, typeExt...)
	}

	return fileExtensions, nil
}

// processFile processes a single file by parsing it, moving it to the output directory if necessary
// and importing the parsed tunes into the database if necessary.
func (fp *FileProcessor) processFile(
	pfc *ProcessFileContext,
	opts *Options,
	pfo *ProcessFilesOptions,
) error {
	tunes, err := fp.parseProcessFile(pfc, opts)
	if errors.Is(err, common.ErrSkipped) {
		return nil
	}
	if err != nil {
		return err
	}

	err = fp.moveProcessFile(pfc.Filepath, opts, pfo)
	if err != nil {
		return err
	}

	if opts.Verbose {
		log.Info().Msgf("(%d/%d) successfully parsed %d tunes from file %s",
			pfc.CurrentFileNr,
			pfc.AllFilesCnt,
			len(tunes),
			pfc.Filepath,
		)
	}

	if !pfo.ImportToDb {
		return nil
	}

	err = fp.importTunesFromFile(tunes, pfc, opts)
	if errors.Is(err, common.ErrSkipped) {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

func (fp *FileProcessor) parseProcessFile(
	pfc *ProcessFileContext,
	opts *Options,
) ([]*messages.ParsedTune, error) {
	tunes, err := fp.parseFile(pfc)
	if err != nil {
		if opts.SkipFailedFiles {
			log.Error().Err(err).Msgf("failed parsing file %s", pfc.Filepath)
			return nil, common.ErrSkipped
		}

		return nil, fmt.Errorf("failed parsing file %s: %v", pfc.Filepath, err)
	}

	return tunes, nil
}

func (fp *FileProcessor) parseFile(
	pfc *ProcessFileContext,
) ([]*messages.ParsedTune, error) {
	fExt := filepath.Ext(pfc.Filepath)
	if fExt == "" {
		return nil, fmt.Errorf(
			"import file %s does not have an extension",
			pfc.Filepath,
		)
	}

	filePlugin, err := fp.pl.PluginForFileExtension(fExt)
	if err != nil {
		return nil, fmt.Errorf("failed getting plugin for file %s with extension %s",
			pfc.Filepath, fExt)
	}

	fileData, err := afero.ReadFile(fp.afs, pfc.Filepath)
	if err != nil {
		return nil, fmt.Errorf("failed reading file %s", pfc.Filepath)
	}

	parsedTunes, err := filePlugin.Parse(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed parsing file %s", pfc.Filepath)
	}

	return parsedTunes, nil
}

// importTunesFromFile imports the given tunes into the database.
// If the import fails, it will return an error if the skip failed files flag is not set.
// If the skip failed files flag is set, it will log the error and return common.ErrSkipped error.
func (fp *FileProcessor) importTunesFromFile(
	parsedTunes []*messages.ParsedTune,
	pfc *ProcessFileContext,
	opts *Options,
) error {
	fInfo, err := fp.fileInfoForImportFile(pfc.Filepath)
	if err != nil {
		if opts.SkipFailedFiles {
			log.Error().Err(err).Msgf("failed getting file info for file %s", pfc.Filepath)
			return nil
		}

		return fmt.Errorf("failed getting file info for file %s: %v", pfc.Filepath, err)
	}

	_, _, err = fp.ds.ImportTunes(parsedTunes, fInfo)
	if err != nil {
		if opts.SkipFailedFiles {
			log.Error().Err(err).Msgf("failed importing parsedTunes from file %s", pfc.Filepath)
			return common.ErrSkipped
		}

		return fmt.Errorf("failed importing parsedTunes from file %s: %v", pfc.Filepath, err)
	}

	return nil
}

func (fp *FileProcessor) fileInfoForImportFile(
	fPath string,
) (*common.ImportFileInfo, error) {
	fExt := filepath.Ext(fPath)
	if fExt == "" {
		return nil, fmt.Errorf("import file %s does not have an extension", fPath)
	}

	fFormat, err := fp.pl.FileFormatForFileExtension(fExt)
	if err != nil {
		return nil, fmt.Errorf("failed getting file format for file %s with extension %s",
			fPath, fExt)
	}

	fInfo, err := common.NewImportFileInfoFromLocalFile(fp.afs, fPath, fFormat)
	if err != nil {
		return nil, fmt.Errorf("failed creating import file info")
	}

	return fInfo, nil
}

// getFilesFromDir returns the absolute paths from all files passed in the paths slice.
// If the path is a directory, it will return all files from that directory.
// If the path is a file and has the correct extension, it will return the file.
// argVerbose is used to print log messages.
// argRecursive is used to get files recursively from the given directories.
func getAllFilesFromPaths(
	afs afero.Fs,
	paths []string,
	o *GetFilesOptions,
) ([]string, error) {
	if len(o.FileExtensions) == 0 {
		return nil, fmt.Errorf("no file extensions given")
	}

	var allFiles []string
	for _, p := range paths {
		fp, err := getFilesFromPath(afs, p, o)
		if err != nil {
			return nil, err
		}
		allFiles = append(allFiles, fp...)
	}

	return allFiles, nil
}

func getFilesFromPath(
	afs afero.Fs,
	path string,
	o *GetFilesOptions,
) ([]string, error) {
	var allFiles []string
	stat, err := afs.Stat(path)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		if isValidFile(stat, o.FileExtensions) {
			allFiles = append(allFiles, path)
		}

		return allFiles, nil
	}

	// if path is a directory, add all files with valid extension to the list
	var dirFiles []string
	if o.Recursive {
		dirFiles, err = getFilesRecursively(afs, path, o.FileExtensions)
	} else {
		dirFiles, err = getFilesFromDir(afs, path, o.FileExtensions)
	}
	if err != nil {
		return nil, err
	}
	allFiles = append(allFiles, dirFiles...)

	return allFiles, nil
}

func isValidFile(stat os.FileInfo, fExt []string) bool {
	return !stat.IsDir() && fileHasValidExtension(stat.Name(), fExt)
}

// getFileFromDir returns the absolute paths from all files in the given directory.
func getFilesFromDir(
	afs afero.Fs,
	path string,
	fExt []string,
) ([]string, error) {
	var filePaths []string
	files, err := afero.ReadDir(afs, path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if fileHasValidExtension(file.Name(), fExt) {
			filePaths = append(filePaths, filepath.Join(path, file.Name()))
		}
	}

	return filePaths, nil
}

// getFilesRecursively returns the absolute paths from all files in the given
// directory and all subdirectories.
func getFilesRecursively(
	afs afero.Fs,
	path string,
	fExt []string,
) ([]string, error) {
	var filePaths []string

	err := afero.Walk(afs, path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if fileHasValidExtension(path, fExt) {
			filePaths = append(filePaths, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return filePaths, nil
}

func fileHasValidExtension(
	path string,
	fExt []string,
) bool {
	ext := filepath.Ext(path)
	return slices.Contains(fExt, ext)
}

// moveProcessFile moves the file to the output directory if the MoveToOutputDir option is set.
func (fp *FileProcessor) moveProcessFile(
	filePath string,
	opts *Options,
	pfo *ProcessFilesOptions,
) error {
	if pfo.MoveToOutputDir {
		err := moveFileToDir(fp.afs, filePath, opts.OutputDir)
		if err != nil {
			log.Error().Err(err).Msgf("failed moving file %s to dir %s", filePath, opts.OutputDir)
			return err
		}
	}

	return nil
}

func moveFileToDir(
	afs afero.Fs,
	file string,
	dir string,
) error {
	if dir == "" {
		return fmt.Errorf("no output directory given")
	}

	fileName := filepath.Base(file)
	newPath := filepath.Join(dir, fileName)
	return afs.Rename(file, newPath)
}

func NewFileProcessor(
	afs afero.Fs,
	pl interfaces.PluginLoader,
	ds interfaces.DataService,
) *FileProcessor {
	return &FileProcessor{
		afs: afs,
		pl:  pl,
		ds:  ds,
	}
}
