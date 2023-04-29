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
	if err := updateTune.Validate(); err != nil {
		return nil, err
	}

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

	for _, apiSet := range apiSets {
		if err := d.setTunesInApiSet(apiSet); err != nil {
			return nil, err
		}
	}

	return apiSets, nil
}

func (d *dbService) CreateMusicSet(musicSet apimodel.CreateSet) (*apimodel.MusicSet, error) {
	dbSet := model.MusicSet{}
	if err := copier.Copy(&dbSet, &musicSet); err != nil {
		return &apimodel.MusicSet{}, fmt.Errorf("could not create db object")
	}

	newTunes, err := d.dbTunesFromIds(musicSet.Tunes)
	if err != nil {
		return nil, err
	}

	err = d.db.Transaction(func(tx *gorm.DB) error {
		if err = d.db.Create(&dbSet).Error; err != nil {
			return err
		}

		if err := d.assignMusicSetTunes(dbSet.ID, musicSet.Tunes); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	dbSet.Tunes = newTunes
	apiSet := &apimodel.MusicSet{}
	if err := copier.Copy(apiSet, dbSet); err != nil {
		return nil, err
	}

	return apiSet, nil
}

func (d *dbService) GetMusicSet(id uint64) (*apimodel.MusicSet, error) {
	var set = &model.MusicSet{}
	if err := d.db.First(set, id).Error; err != nil {
		return &apimodel.MusicSet{}, common.NotFound
	}

	apiSet := &apimodel.MusicSet{}
	if err := copier.Copy(apiSet, set); err != nil {
		return &apimodel.MusicSet{}, err
	}

	if err := d.setTunesInApiSet(apiSet); err != nil {
		return &apimodel.MusicSet{}, err
	}

	return apiSet, nil
}

func (d *dbService) setTunesInApiSet(apiSet *apimodel.MusicSet) error {
	var setTunes []model.Tune
	err := d.db.Joins("JOIN music_set_tunes mst on tunes.id = mst.tune_id").
		Where("mst.music_set_id=?", apiSet.ID).
		Order("mst.\"order\"").
		Find(&setTunes).Error
	if err != nil {
		return fmt.Errorf("failed getting music set tunes: %s", err.Error())
	}

	if err := copier.Copy(&apiSet.Tunes, &setTunes); err != nil {
		return err
	}

	return nil
}

func (d *dbService) UpdateMusicSet(id uint64, updateSet apimodel.UpdateSet) (*apimodel.MusicSet, error) {
	if err := updateSet.Validate(); err != nil {
		return nil, err
	}

	var dbSet = &model.MusicSet{}
	if err := d.db.Preload("Tunes").First(dbSet, id).Error; err != nil {
		return nil, common.NotFound
	}

	var updateVals = map[string]interface{}{}
	if err := mapstructure.Decode(&updateSet, &updateVals); err != nil {
		return nil, err
	}

	newTunes, err := d.dbTunesFromIds(updateSet.Tunes)
	if err != nil {
		return nil, err
	}

	err = d.db.Transaction(func(tx *gorm.DB) error {
		if err := d.deleteMusicSetTunes(dbSet); err != nil {
			return err
		}

		if err := d.db.Model(dbSet).Updates(updateVals).Error; err != nil {
			return err
		}

		if err := d.assignMusicSetTunes(dbSet.ID, updateSet.Tunes); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	dbSet.Tunes = newTunes
	apiSet := &apimodel.MusicSet{}
	if err := copier.Copy(apiSet, dbSet); err != nil {
		return nil, err
	}

	return apiSet, nil
}

func (d *dbService) DeleteMusicSet(id uint64) error {
	var set = &model.MusicSet{}
	if err := d.db.Preload("Tunes").First(set, id).Error; err != nil {
		return common.NotFound
	}

	err := d.db.Transaction(func(tx *gorm.DB) error {
		if err := d.deleteMusicSetTunes(set); err != nil {
			return err
		}

		if err := d.db.Delete(set).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (d *dbService) AssignTunesToMusicSet(
	setId uint64,
	tuneIds []uint64,
) (*apimodel.MusicSet, error) {
	set := &model.MusicSet{}
	if err := d.db.Preload("Tunes").First(set, setId).Error; err != nil {
		return nil, common.NotFound
	}

	newTunes, err := d.dbTunesFromIds(tuneIds)
	if err != nil {
		return nil, err
	}

	// delete old music set -> tune relations and create new ones
	err = d.db.Transaction(func(tx *gorm.DB) error {
		if err := d.deleteMusicSetTunes(set); err != nil {
			return err
		}

		if err := d.assignMusicSetTunes(set.ID, tuneIds); err != nil {
			return err
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

// dbTunesFromIds returns the database tune objects in the same order as the
// given tuneIds. If there is an id that belongs to a non existing tune,
// an error will be returned.
func (d *dbService) dbTunesFromIds(tuneIds []uint64) ([]model.Tune, error) {
	if len(tuneIds) == 0 {
		return nil, nil
	}

	var dbTunes []model.Tune
	if err := d.db.Where("id IN (?)", tuneIds).Find(&dbTunes).Error; err != nil {
		return nil, err
	}

	if len(dbTunes) != len(tuneIds) {
		return nil, fmt.Errorf("not all tune IDs are from valid tunes")
	}

	return tunesOrderedByIds(dbTunes, tuneIds), nil
}

func (d *dbService) deleteMusicSetTunes(set *model.MusicSet) error {
	for _, tune := range set.Tunes {
		if err := d.db.Delete(&model.MusicSetTunes{
			MusicSetID: set.ID,
			TuneID:     tune.ID,
		}).Error; err != nil {
			return err
		}
	}

	// set tunes to nil to reflect the database state
	set.Tunes = nil

	return nil
}

func (d *dbService) assignMusicSetTunes(setId uint64, tuneIds []uint64) error {
	for i, tuneId := range tuneIds {
		setTune := &model.MusicSetTunes{
			MusicSetID: setId,
			TuneID:     tuneId,
			Order:      uint(i + 1),
		}
		if err := d.db.Create(setTune).Error; err != nil {
			return err
		}
	}

	return nil
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
