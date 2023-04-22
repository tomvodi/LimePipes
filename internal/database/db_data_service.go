package database

import (
	"banduslib/internal/api/apimodel"
	"banduslib/internal/common"
	"banduslib/internal/database/model"
	"banduslib/internal/interfaces"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type dbService struct {
	db *gorm.DB
}

func (d *dbService) Tunes() ([]*apimodel.Tune, error) {
	var tunes []model.Tune
	if err := d.db.Find(&tunes).Error; err != nil {
		return nil, err
	}

	var apiTunes []*apimodel.Tune
	if err := copier.Copy(&apiTunes, tunes); err != nil {
		return nil, err
	}

	return apiTunes, nil
}

func (d *dbService) CreateTune(tune apimodel.CreateTune) (*apimodel.Tune, error) {
	dbTune := model.Tune{}
	err := copier.Copy(&dbTune, &tune)
	if err != nil {
		return &apimodel.Tune{}, fmt.Errorf("could not create db object")
	}

	if err = d.db.Create(&dbTune).Error; err != nil {
		return &apimodel.Tune{}, err
	}

	apiTune := &apimodel.Tune{}
	if err := copier.Copy(apiTune, dbTune); err != nil {
		return nil, err
	}

	return apiTune, nil
}

func (d *dbService) GetTune(id uint64) (*apimodel.Tune, error) {
	var tune = &model.Tune{}
	if err := d.db.Preload("Sets").First(tune, id).Error; err != nil {
		return &apimodel.Tune{}, common.NotFound
	}

	apiTune := &apimodel.Tune{}
	err := copier.Copy(apiTune, tune)
	if err != nil {
		return &apimodel.Tune{}, err
	}

	return apiTune, nil
}

func (d *dbService) UpdateTune(id uint64, updateTune apimodel.UpdateTune) (*apimodel.Tune, error) {
	var tune = &model.Tune{}
	if err := d.db.First(tune, id).Error; err != nil {
		return nil, common.NotFound
	}

	var updateVals = map[string]interface{}{}
	if err := mapstructure.Decode(&updateTune, &updateVals); err != nil {
		return nil, err
	}

	if err := d.db.Model(tune).Updates(updateVals).Error; err != nil {
		return nil, err
	}

	apiTune := &apimodel.Tune{}
	if err := copier.Copy(apiTune, tune); err != nil {
		return nil, err
	}

	return apiTune, nil
}

func (d *dbService) DeleteTune(id uint64) error {
	var tune = &model.Tune{}
	if err := d.db.First(tune, id).Error; err != nil {
		return common.NotFound
	}

	if err := d.db.Delete(tune).Error; err != nil {
		return err
	}

	return nil
}

func (d *dbService) MusicSets() ([]*apimodel.MusicSet, error) {
	var sets []model.MusicSet
	if err := d.db.Find(&sets).Error; err != nil {
		return nil, err
	}

	var apiSets []*apimodel.MusicSet
	if err := copier.Copy(&apiSets, sets); err != nil {
		return nil, err
	}

	return apiSets, nil
}

func (d *dbService) CreateMusicSet(tune apimodel.CreateSet) (*apimodel.MusicSet, error) {
	dbSet := model.MusicSet{}
	err := copier.Copy(&dbSet, &tune)
	if err != nil {
		return &apimodel.MusicSet{}, fmt.Errorf("could not create db object")
	}

	if err = d.db.Create(&dbSet).Error; err != nil {
		return &apimodel.MusicSet{}, err
	}

	apiTune := &apimodel.MusicSet{}
	if err := copier.Copy(apiTune, dbSet); err != nil {
		return nil, err
	}

	return apiTune, nil
}

func (d *dbService) GetMusicSet(id uint64) (*apimodel.MusicSet, error) {
	var set = &model.MusicSet{}
	if err := d.db.First(set, id).Error; err != nil {
		return &apimodel.MusicSet{}, common.NotFound
	}

	var setTunes []model.Tune
	err := d.db.Joins("JOIN music_set_tunes mst on tunes.id = mst.tune_id").
		Where("mst.music_set_id=?", set.ID).
		Order("mst.\"order\"").
		Find(&setTunes).Error
	if err != nil {
		return &apimodel.MusicSet{},
			fmt.Errorf("failed getting music set tunes: %s", err.Error())
	}

	apiSet := &apimodel.MusicSet{}
	if err := copier.Copy(apiSet, set); err != nil {
		return &apimodel.MusicSet{}, err
	}
	if err := copier.Copy(&apiSet.Tunes, &setTunes); err != nil {
		return &apimodel.MusicSet{}, err
	}

	return apiSet, nil
}

func (d *dbService) UpdateMusicSet(id uint64, updateSet apimodel.UpdateSet) (*apimodel.MusicSet, error) {
	var set = &model.MusicSet{}
	if err := d.db.First(set, id).Error; err != nil {
		return nil, common.NotFound
	}

	var updateVals = map[string]interface{}{}
	if err := mapstructure.Decode(&updateSet, &updateVals); err != nil {
		return nil, err
	}

	if err := d.db.Model(set).Updates(updateVals).Error; err != nil {
		return nil, err
	}

	apiSet := &apimodel.MusicSet{}
	if err := copier.Copy(apiSet, set); err != nil {
		return nil, err
	}

	return apiSet, nil
}

func (d *dbService) DeleteMusicSet(id uint64) error {
	var set = &model.MusicSet{}
	if err := d.db.First(set, id).Error; err != nil {
		return common.NotFound
	}

	if err := d.db.Delete(set).Error; err != nil {
		return err
	}

	return nil
}

func (d *dbService) AssignTunesToMusicSet(
	setId uint64,
	tuneIds []uint64,
) (*apimodel.MusicSet, error) {
	set := &model.MusicSet{}
	if err := d.db.Preload("Tunes").First(set, setId).Error; err != nil {
		return nil, common.NotFound
	}

	var newTunes []model.Tune
	if err := d.db.Where("id IN (?)", tuneIds).Find(&newTunes).Error; err != nil {
		return nil, err
	}

	if len(newTunes) != len(tuneIds) {
		return nil, fmt.Errorf("not all tune IDs are from valid tunes")
	}
	newTunes = tunesOrderedByIds(newTunes, tuneIds)

	// delete old music set -> tune relations and create new ones
	err := d.db.Transaction(func(tx *gorm.DB) error {
		for _, tune := range set.Tunes {
			if err := d.db.Delete(&model.MusicSetTunes{
				MusicSetID: setId,
				TuneID:     tune.ID,
			}).Error; err != nil {
				return err
			}
		}

		for i, tune := range newTunes {
			setTune := &model.MusicSetTunes{
				MusicSetID: setId,
				TuneID:     tune.ID,
				Order:      uint(i + 1),
			}
			if err := d.db.Create(setTune).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Return new set with tunes in order
	set.Tunes = newTunes

	apiSet := &apimodel.MusicSet{}
	if err := copier.Copy(apiSet, set); err != nil {
		return nil, err
	}

	return apiSet, nil
}

func tunesOrderedByIds(tunes []model.Tune, tuneIds []uint64) []model.Tune {
	var orderedTunes = make([]model.Tune, len(tunes))
	for i, id := range tuneIds {
		for _, tune := range tunes {
			if tune.ID == id {
				orderedTunes[i] = tune
				break
			}
		}
	}
	return orderedTunes
}

func NewDbDataService(db *gorm.DB) interfaces.DataService {
	return &dbService{
		db: db,
	}
}
