package common

import "os"

type ImportFileInfo struct {
	OriginalPath string
	Name         string
	Hash         string
	Data         []byte
}

func NewImportFileInfoFromLocalFile(originalPath string) (*ImportFileInfo, error) {
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
		Hash:         fHash,
		Data:         fileData,
	}

	return fInfo, nil
}

func NewImportFileInfo(fileName string, fileData []byte) (*ImportFileInfo, error) {
	fHash, err := HashFromData(fileData)
	if err != nil {
		return nil, err
	}

	fInfo := &ImportFileInfo{
		OriginalPath: fileName,
		Name:         FilenameFromPath(fileName),
		Hash:         fHash,
		Data:         fileData,
	}

	return fInfo, nil
}
