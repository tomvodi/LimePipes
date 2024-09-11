package interfaces

import (
	"github.com/google/uuid"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"github.com/tomvodi/limepipes/internal/apigen/apimodel"
	"github.com/tomvodi/limepipes/internal/common"
	"github.com/tomvodi/limepipes/internal/database/model"
)

type DataService interface {
	Tunes() ([]*apimodel.Tune, error)
	CreateTune(tune apimodel.CreateTune, importFile *model.ImportFile) (*apimodel.Tune, error)
	GetTune(id uuid.UUID) (*apimodel.Tune, error)
	UpdateTune(id uuid.UUID, tune apimodel.UpdateTune) (*apimodel.Tune, error)
	DeleteTune(id uuid.UUID) error

	AddFileToTune(tuneID uuid.UUID, tFile *model.TuneFile) error
	DeleteFileFromTune(tuneID uuid.UUID, fType fileformat.Format) error
	GetTuneFile(tuneID uuid.UUID, fType fileformat.Format) (*model.TuneFile, error)
	GetTuneFiles(tuneID uuid.UUID) ([]*model.TuneFile, error)

	MusicSets() ([]*apimodel.MusicSet, error)
	CreateMusicSet(tune apimodel.CreateSet, importFile *model.ImportFile) (*apimodel.MusicSet, error)
	GetMusicSet(id uuid.UUID) (*apimodel.MusicSet, error)
	UpdateMusicSet(id uuid.UUID, tune apimodel.UpdateSet) (*apimodel.MusicSet, error)
	DeleteMusicSet(id uuid.UUID) error

	AssignTunesToMusicSet(setID uuid.UUID, tuneIDs []uuid.UUID) (*apimodel.MusicSet, error)

	GetImportFileByHash(fHash string) (*model.ImportFile, error)
	ImportTunes(
		tunes []*messages.ImportedTune,
		fileInfo *common.ImportFileInfo,
	) ([]*apimodel.ImportTune, *apimodel.BasicMusicSet, error)
}
