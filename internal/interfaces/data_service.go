package interfaces

import (
	"github.com/google/uuid"
	"github.com/tomvodi/limepipes/internal/api_gen/apimodel"
	"github.com/tomvodi/limepipes/internal/common"
	"github.com/tomvodi/limepipes/internal/common/music_model"
	"github.com/tomvodi/limepipes/internal/database/model"
	"github.com/tomvodi/limepipes/internal/database/model/file_type"
)

type DataService interface {
	Tunes() ([]*apimodel.Tune, error)
	CreateTune(tune apimodel.CreateTune, importFile *model.ImportFile) (*apimodel.Tune, error)
	GetTune(id uuid.UUID) (*apimodel.Tune, error)
	UpdateTune(id uuid.UUID, tune apimodel.UpdateTune) (*apimodel.Tune, error)
	DeleteTune(id uuid.UUID) error

	AddFileToTune(tuneId uuid.UUID, tFile *model.TuneFile) error
	DeleteFileFromTune(tuneId uuid.UUID, fType file_type.Type) error
	GetTuneFile(tuneId uuid.UUID, fType file_type.Type) (*model.TuneFile, error)
	GetTuneFiles(tuneId uuid.UUID) ([]*model.TuneFile, error)

	MusicSets() ([]*apimodel.MusicSet, error)
	CreateMusicSet(tune apimodel.CreateSet, importFile *model.ImportFile) (*apimodel.MusicSet, error)
	GetMusicSet(id uuid.UUID) (*apimodel.MusicSet, error)
	UpdateMusicSet(id uuid.UUID, tune apimodel.UpdateSet) (*apimodel.MusicSet, error)
	DeleteMusicSet(id uuid.UUID) error

	AssignTunesToMusicSet(setId uuid.UUID, tuneIds []uuid.UUID) (*apimodel.MusicSet, error)

	ImportMusicModel(
		muMo music_model.MusicModel,
		fileInfo *common.ImportFileInfo,
		bwwFileData *common.BwwFileTuneData,
	) ([]*apimodel.ImportTune, error)
}
