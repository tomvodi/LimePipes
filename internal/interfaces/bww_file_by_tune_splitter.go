package interfaces

import "github.com/tomvodi/limepipes/internal/common"

//go:generate mockgen -source bww_file_by_tune_splitter.go -destination ./mocks/mock_bww_file_by_tune_splitter.go

type BwwFileByTuneSplitter interface {
	SplitFileData(data []byte) (*common.BwwFileTuneData, error)
}
