package model

type ImportFile struct {
	BaseModel
	Name         string
	OriginalPath string
	Hash         string
	Data         []byte // file content from original file
}
