package database

import (
	"banduslib/internal/api/apimodel"
	"banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/common/music_model/import_message"
	"banduslib/internal/database/model"
	"banduslib/internal/database/model/file_type"
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

	var distinctTuneIds []uint64
	for _, id := range tuneIds {
		inDistinct := false
		for _, distTuneId := range distinctTuneIds {
			if id == distTuneId {
				inDistinct = true
			}
		}
		if !inDistinct {
			distinctTuneIds = append(distinctTuneIds, id)
		}
	}

	var dbTunes []model.Tune
	if err := d.db.Where("id IN (?)", distinctTuneIds).Find(&dbTunes).Error; err != nil {
		return nil, err
	}

	if len(dbTunes) != len(distinctTuneIds) {
		return nil, fmt.Errorf("not all tune IDs are from valid tunes")
	}

	return tunesOrderedByIds(dbTunes, tuneIds), nil
}

func (d *dbService) deleteMusicSetTunes(set *model.MusicSet) error {
	err := d.db.Where(&model.MusicSetTunes{
		MusicSetID: set.ID,
	}).Delete(&model.MusicSetTunes{}).Error

	if err != nil {
		return err
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
	var orderedTunes = make([]model.Tune, len(tuneIds))
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

func (d *dbService) ImportMusicModel(
	muMo music_model.MusicModel,
	filename string,
	bwwFileData *common.BwwFileTuneData,
) ([]*apimodel.ImportTune, error) {
	var apiTunes []*apimodel.ImportTune

	err := d.db.Transaction(func(tx *gorm.DB) error {
		for _, tune := range muMo {
			timeSigStr := ""
			timeSig := tune.FirstTimeSignature()
			if timeSig != nil {
				timeSigStr = timeSig.String()
			}
			alreadyImportedTune := apiTuneWithTitle(tune.Title, apiTunes)
			if alreadyImportedTune != nil {
				apiTunes = append(apiTunes, alreadyImportedTune)
				continue
			}

			createTune := apimodel.CreateTune{
				Title:    tune.Title,
				Type:     tune.Type,
				TimeSig:  timeSigStr,
				Composer: tune.Composer,
			}
			apiTune, err := d.CreateTune(createTune)
			if err != nil {
				return err
			}
			tuneFile, err := model.TuneFileFromTune(tune)
			if err != nil {
				return err
			}

			if err = d.AddFileToTune(apiTune.ID, tuneFile); err != nil {
				return err
			}

			if bwwFileData != nil {
				if !bwwFileData.HasDataForTune(tune.Title) {
					return fmt.Errorf("no bww tune file data for tune %s", tune.Title)
				}

				tuneFile = &model.TuneFile{
					Type: file_type.Bww,
					Data: bwwFileData.DataForTune(tune.Title),
				}
				if err = d.AddFileToTune(apiTune.ID, tuneFile); err != nil {
					return err
				}
			}

			importTune := &apimodel.ImportTune{}
			err = copier.Copy(&importTune, apiTune)
			if err != nil {
				return err
			}
			setMessagesToApiTune(importTune, tune)
			apiTunes = append(apiTunes, importTune)
		}

		if len(muMo) > 1 {
			var tuneIds []uint64
			for _, tune := range apiTunes {
				tuneIds = append(tuneIds, tune.ID)
			}
			createSet := apimodel.CreateSet{
				Title:       filename,
				Description: "imported from bww file",
				Creator:     "",
				Tunes:       tuneIds,
			}
			apiSet, err := d.CreateMusicSet(createSet)
			if err != nil {
				return fmt.Errorf("failed creating set for file %s: %s", filename, err.Error())
			}

			basicSet := &apimodel.BasicMusicSet{}
			err = copier.Copy(basicSet, apiSet)
			if err != nil {
				return fmt.Errorf("failed creating basic music set from music set")
			}

			for _, tune := range apiTunes {
				tune.Set = basicSet
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return apiTunes, nil
}

func apiTuneWithTitle(
	title string,
	apiTunes []*apimodel.ImportTune,
) *apimodel.ImportTune {
	if apiTunes == nil {
		return nil
	}

	for _, tune := range apiTunes {
		if tune.Title == title {
			return tune
		}
	}
	return nil
}

func setMessagesToApiTune(apiTune *apimodel.ImportTune, tune *music_model.Tune) {
	messages := tune.ImportMessages()
	for _, message := range messages {
		switch message.Type {
		case import_message.Error:
			apiTune.Errors = append(apiTune.Errors, message.Text)
		case import_message.Warning:
			apiTune.Warnings = append(apiTune.Warnings, message.Text)
		case import_message.Info:
			apiTune.Infos = append(apiTune.Infos, message.Text)
		}
	}
}

func (d *dbService) GetTuneFile(tuneId uint64, fType file_type.Type) (*model.TuneFile, error) {
	tuneFile := &model.TuneFile{
		TuneID: tuneId,
		Type:   fType,
	}

	if err := d.db.First(tuneFile).Error; err != nil {
		return nil, common.NotFound
	}

	return tuneFile, nil
}

func (d *dbService) GetTuneFiles(tuneId uint64) ([]*model.TuneFile, error) {
	var tune = &model.Tune{}
	if err := d.db.First(tune, tuneId).Error; err != nil {
		return nil, common.NotFound
	}

	var tuneFiles []*model.TuneFile
	if err := d.db.Where("tune_id = ?", tuneId).Find(&tuneFiles).Error; err != nil {
		return nil, err
	}

	return tuneFiles, nil
}

func (d *dbService) AddFileToTune(tuneId uint64, tFile *model.TuneFile) error {
	var tune = &model.Tune{}
	if err := d.db.First(tune, tuneId).Error; err != nil {
		return common.NotFound
	}

	tFile.TuneID = tuneId

	if err := d.db.Create(tFile).Error; err != nil {
		return err
	}

	return nil
}

func (d *dbService) DeleteFileFromTune(tuneId uint64, fType file_type.Type) error {
	tuneFile := &model.TuneFile{
		TuneID: tuneId,
		Type:   fType,
	}
	if err := d.db.Delete(tuneFile).Error; err != nil {
		return err
	}

	return nil
}

func NewDbDataService(db *gorm.DB) interfaces.DataService {
	return &dbService{
		db: db,
	}
}
