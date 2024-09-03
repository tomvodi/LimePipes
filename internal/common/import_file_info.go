package common

import (
	"github.com/spf13/afero"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
)

type ImportFileInfo struct {
	OriginalPath string
	FileFormat   fileformat.Format
	Name         string
	Hash         string
	Data         []byte
}

func NewImportFileInfoFromLocalFile(
	afs afero.Fs,
	originalPath string,
	fFormat fileformat.Format,
) (*ImportFileInfo, error) {
	fHash, err := HashFromFile(afs, originalPath)
	if err != nil {
		return nil, err
	}

	fileData, err := afero.ReadFile(afs, originalPath)
	if err != nil {
		return nil, err
	}

	fInfo := &ImportFileInfo{
		OriginalPath: originalPath,
		Name:         FilenameFromPath(originalPath),
		FileFormat:   fFormat,
		Hash:         fHash,
		Data:         fileData,
	}

	return fInfo, nil
}

func NewImportFileInfo(
	fileName string,
	fFormat fileformat.Format,
	fileData []byte,
) (*ImportFileInfo, error) {
	fHash, err := HashFromData(fileData)
	if err != nil {
		return nil, err
	}

	fInfo := &ImportFileInfo{
		OriginalPath: fileName,
		Name:         FilenameFromPath(fileName),
		FileFormat:   fFormat,
		Hash:         fHash,
		Data:         fileData,
	}

	return fInfo, nil
}
