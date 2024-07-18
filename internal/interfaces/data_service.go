package interfaces

import (
	"github.com/tomvodi/limepipes/internal/api_gen/apimodel"
	"github.com/tomvodi/limepipes/internal/common"
	"github.com/tomvodi/limepipes/internal/common/music_model"
	"github.com/tomvodi/limepipes/internal/database/model"
	"github.com/tomvodi/limepipes/internal/database/model/file_type"
)

//go:generate mockgen -source data_service.go -destination ./mocks/mock_data_service.go

type DataService interface {
	Tunes() ([]*apimodel.Tune, error)
	CreateTune(tune apimodel.CreateTune, importFile *model.ImportFile) (*apimodel.Tune, error)
	GetTune(id int64) (*apimodel.Tune, error)
	UpdateTune(id int64, tune apimodel.UpdateTune) (*apimodel.Tune, error)
	DeleteTune(id int64) error

	AddFileToTune(tuneId int64, tFile *model.TuneFile) error
	DeleteFileFromTune(tuneId int64, fType file_type.Type) error
	GetTuneFile(tuneId int64, fType file_type.Type) (*model.TuneFile, error)
	GetTuneFiles(tuneId int64) ([]*model.TuneFile, error)

	MusicSets() ([]*apimodel.MusicSet, error)
	CreateMusicSet(tune apimodel.CreateSet, importFile *model.ImportFile) (*apimodel.MusicSet, error)
	GetMusicSet(id int64) (*apimodel.MusicSet, error)
	UpdateMusicSet(id int64, tune apimodel.UpdateSet) (*apimodel.MusicSet, error)
	DeleteMusicSet(id int64) error

	AssignTunesToMusicSet(setId int64, tuneIds []int64) (*apimodel.MusicSet, error)

	ImportMusicModel(
		muMo music_model.MusicModel,
		fileInfo *common.ImportFileInfo,
		bwwFileData *common.BwwFileTuneData,
	) ([]*apimodel.ImportTune, error)
}
