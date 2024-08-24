package common

import (
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/file_type"
	"os"
)

type ImportFileInfo struct {
	OriginalPath string
	FileType     file_type.Type
	Name         string
	Hash         string
	Data         []byte
}

func NewImportFileInfoFromLocalFile(
	originalPath string,
	fType file_type.Type,
) (*ImportFileInfo, error) {
	fHash, err := HashFromFile(originalPath)
	if err != nil {
		return nil, err
	}

	fileData, err := os.ReadFile(originalPath)
	if err != nil {
		return nil, err
	}

	fInfo := &ImportFileInfo{
		OriginalPath: originalPath,
		Name:         FilenameFromPath(originalPath),
		FileType:     fType,
		Hash:         fHash,
		Data:         fileData,
	}

	return fInfo, nil
}

func NewImportFileInfo(fileName string, fType file_type.Type, fileData []byte) (*ImportFileInfo, error) {
	fHash, err := HashFromData(fileData)
	if err != nil {
		return nil, err
	}

	fInfo := &ImportFileInfo{
		OriginalPath: fileName,
		Name:         FilenameFromPath(fileName),
		FileType:     fType,
		Hash:         fHash,
		Data:         fileData,
	}

	return fInfo, nil
}
