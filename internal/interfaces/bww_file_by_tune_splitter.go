package interfaces

import "github.com/tomvodi/limepipes/internal/common"

type BwwFileByTuneSplitter interface {
	SplitFileData(data []byte) (*common.BwwFileTuneData, error)
}
