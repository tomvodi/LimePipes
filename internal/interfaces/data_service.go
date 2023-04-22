package interfaces

import "banduslib/internal/api/apimodel"

type DataService interface {
	Tunes() ([]*apimodel.Tune, error)
	CreateTune(tune apimodel.CreateTune) (*apimodel.Tune, error)
	GetTune(id uint64) (*apimodel.Tune, error)
	UpdateTune(id uint64, tune apimodel.UpdateTune) (*apimodel.Tune, error)
	DeleteTune(id uint64) error

	MusicSets() ([]*apimodel.MusicSet, error)
	CreateMusicSet(tune apimodel.CreateSet) (*apimodel.MusicSet, error)
	GetMusicSet(id uint64) (*apimodel.MusicSet, error)
	UpdateMusicSet(id uint64, tune apimodel.UpdateSet) (*apimodel.MusicSet, error)
	DeleteMusicSet(id uint64) error

	AssignTunesToMusicSet(setId uint64, tuneIds []uint64) (*apimodel.MusicSet, error)
}
