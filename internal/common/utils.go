package common

import (
	"path/filepath"
	"strings"
)

// FilenameFromPath returns for a given filename with extension or filepath
// only the filename without an extension and without the path
func FilenameFromPath(file string) string {
	onlyFile := filepath.Base(file)
	return strings.TrimSuffix(onlyFile, filepath.Ext(onlyFile))
}
