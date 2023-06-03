package common

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FilenameFromPath returns for a given filename with extension or filepath
// only the filename without an extension and without the path
func FilenameFromPath(file string) string {
	onlyFile := filepath.Base(file)
	return strings.TrimSuffix(onlyFile, filepath.Ext(onlyFile))
}

func HashFromFile(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func HashFromData(data []byte) (string, error) {
	reader := bytes.NewReader(data)

	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
