// Package importtype provides the valid import types for the limepipes-cli.
// These may differ from actual file formats supported by limepipes plugins. Also, the import types
// are lower case strings while the file formats are upper case strings.
package importtype

import (
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
	"golang.org/x/exp/maps"
	"strings"
)

// FileFormatMapping returns a mapping of import type strings used in the cli to file format enums.
func FileFormatMapping() map[string]fileformat.Format {
	return map[string]fileformat.Format{
		FromFileFormat(fileformat.Format_BWW): fileformat.Format_BWW,
	}
}

func AllTypes() (allK []string) {
	allK = append(allK, maps.Keys(FileFormatMapping())...)
	return allK
}

func FromFileFormat(ff fileformat.Format) string {
	return strings.ToLower(ff.String())
}
