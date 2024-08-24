package database

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/measure"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/tune"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/file_type"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"github.com/tomvodi/limepipes/internal/api_gen/apimodel"
	"github.com/tomvodi/limepipes/internal/common"

	"github.com/tomvodi/limepipes/internal/database/model"
	"github.com/tomvodi/limepipes/internal/interfaces"
	"gorm.io/gorm"
	"strings"
)

type dbService struct {
	db        *gorm.DB
	validator interfaces.ApiModelValidator
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

func (d *dbService) CreateTune(
	tune apimodel.CreateTune,
	importFile *model.ImportFile,
) (*apimodel.Tune, error) {
	if strings.TrimSpace(tune.Title) == "" {
		return nil, fmt.Errorf("can't create tune without a title")
	}

	dbTune := model.Tune{}
	if importFile != nil {
		dbTune.ImportFileId = importFile.ID
	}
	var tuneType *model.TuneType
	var err error
	if tune.Type != "" {
		tuneType, err = d.getOrCreateTuneType(tune.Type)
		if err != nil {
			return &apimodel.Tune{}, err
		}
	}

	err = copier.Copy(&dbTune, &tune)
	if err != nil {
		return &apimodel.Tune{}, fmt.Errorf("could not create db object")
	}
	if tuneType != nil {
		dbTune.TuneTypeId = &tuneType.ID
	}

	if err = d.db.Create(&dbTune).Error; err != nil {
		return &apimodel.Tune{}, err
	}

	apiTune := &apimodel.Tune{}
	if err := copier.Copy(apiTune, &dbTune); err != nil {
		return nil, err
	}
	if tuneType != nil {
		apiTune.Type = tuneType.Name
	}

	return apiTune, nil
}

func (d *dbService) GetTune(id uuid.UUID) (*apimodel.Tune, error) {
	var tune = &model.Tune{}
	if err := d.db.
		Preload("Sets").
		Preload("TuneType").
		First(tune, id).Error; err != nil {
		return &apimodel.Tune{}, common.NotFound
	}

	apiTune := &apimodel.Tune{}
	err := copier.Copy(apiTune, tune)
	if err != nil {
		return &apimodel.Tune{}, err
	}
	if tune.TuneType != nil {
		apiTune.Type = tune.TuneType.Name
	}

	return apiTune, nil
}

// getTuneFromDb returns a tune from the database if it is identical
// this means it has the same title, composer, type and arranger
func (d *dbService) getIdenticalTuneFromDb(tune *tune.Tune) (*apimodel.Tune, error) {
	dbTune, err := d.getTuneByTitle(tune.Title)
	if err != nil {
		return nil, err
	}

	if dbTune.Composer == tune.Composer &&
		dbTune.Arranger == tune.Arranger &&
		dbTune.Type == tune.Type {
		return dbTune, nil
	}

	return nil, common.NotFound
}

func (d *dbService) getTuneByTitle(title string) (*apimodel.Tune, error) {
	var tune = &model.Tune{
		Title: title,
	}
	if err := d.db.
		Preload("TuneType").
		Where(tune).
		First(tune).Error; err != nil {
		return nil, common.NotFound
	}

	apiTune := &apimodel.Tune{}
	err := copier.Copy(apiTune, tune)
	if err != nil {
		return nil, err
	}
	if tune.TuneType != nil {
		apiTune.Type = tune.TuneType.Name
	}

	return apiTune, nil
}

func (d *dbService) GetImportFileByHash(fHash string) (*model.ImportFile, error) {
	var importFile = &model.ImportFile{
		Hash: fHash,
	}
	if err := d.db.Where(importFile).First(importFile).Error; err != nil {
		return nil, common.NotFound
	}

	return importFile, nil
}

func (d *dbService) hasImportFile(
	fileInfo *common.ImportFileInfo,
) (bool, error) {
	_, err := d.GetImportFileByHash(fileInfo.Hash)
	if errors.Is(err, common.NotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, err
}

func (d *dbService) createImportFile(importFile *common.ImportFileInfo) (*model.ImportFile, error) {
	dbImportFile := &model.ImportFile{}
	err := copier.Copy(dbImportFile, importFile)
	if err != nil {
		return nil, err
	}

	if err = d.db.Create(&dbImportFile).Error; err != nil {
		return nil, err
	}

	return dbImportFile, nil
}

func (d *dbService) UpdateTune(id uuid.UUID, updateTune apimodel.UpdateTune) (*apimodel.Tune, error) {
	if err := d.validator.ValidateUpdateTune(updateTune); err != nil {
		return nil, err
	}

	tuneType, err := d.getOrCreateTuneType(updateTune.Type)
	if err != nil {
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
	updateVals["TuneTypeId"] = tuneType.ID
	delete(updateVals, "Type")

	if err := d.db.Model(tune).Updates(updateVals).Error; err != nil {
		return nil, err
	}

	apiTune := &apimodel.Tune{}
	if err := copier.Copy(apiTune, tune); err != nil {
		return nil, err
	}
	apiTune.Type = tuneType.Name

	return apiTune, nil
}

func (d *dbService) DeleteTune(id uuid.UUID) error {
	var tune = &model.Tune{}
	if err := d.db.First(tune, id).Error; err != nil {
		return common.NotFound
	}

	if err := d.db.Delete(tune).Error; err != nil {
		return err
	}

	return nil
}

func (d *dbService) getOrCreateTuneType(
	name string,
) (*model.TuneType, error) {
	var tuneType = &model.TuneType{}
	var err error
	tuneType, err = d.getTuneTypeByName(name)
	if errors.Is(err, common.NotFound) {
		return d.createTuneType(name)
	}

	return tuneType, err
}

func (d *dbService) getTuneTypeByName(name string) (*model.TuneType, error) {
	var tuneType = &model.TuneType{}
	if err := d.db.Where("lower(name) = ?", strings.ToLower(name)).
		First(tuneType).Error; err != nil {
		return nil, common.NotFound
	}

	return tuneType, nil
}

func (d *dbService) createTuneType(name string) (*model.TuneType, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("tune type name must have a value")
	}

	dbTuneType := &model.TuneType{
		Name: name,
	}

	if err := d.db.Create(&dbTuneType).Error; err != nil {
		return nil, err
	}

	return dbTuneType, nil
}

func (d *dbService) MusicSets() ([]*apimodel.MusicSet, error) {
	var sets []model.MusicSet
	if err := d.db.Find(&sets).Error; err != nil {
		return nil, err
	}

	apiSets := make([]*apimodel.MusicSet, len(sets))
	for i, set := range sets {
		var err error
		apiSets[i], err = apiSetFromDbSet(&set)
		if err != nil {
			return nil, err
		}
	}

	for _, apiSet := range apiSets {
		if err := d.setTunesInApiSet(apiSet); err != nil {
			return nil, err
		}
	}

	return apiSets, nil
}

func (d *dbService) CreateMusicSet(
	musicSet apimodel.CreateSet,
	importFile *model.ImportFile,
) (*apimodel.MusicSet, error) {
	if strings.TrimSpace(musicSet.Title) == "" {
		return nil, fmt.Errorf("can't create music set without a title")
	}

	dbSet := model.MusicSet{}
	if importFile != nil {
		dbSet.ImportFileId = importFile.ID
	}
	if err := copier.Copy(&dbSet, &musicSet); err != nil {
		return &apimodel.MusicSet{}, fmt.Errorf("could not create db object")
	}
	dbSet.Tunes = nil

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
	apiSet, err := apiSetFromDbSet(&dbSet)
	if err != nil {
		return &apimodel.MusicSet{}, err
	}

	return apiSet, nil
}

func (d *dbService) GetMusicSet(id uuid.UUID) (*apimodel.MusicSet, error) {
	var set = &model.MusicSet{}
	if err := d.db.First(set, id).Error; err != nil {
		return &apimodel.MusicSet{}, common.NotFound
	}

	apiSet, err := apiSetFromDbSet(set)
	if err != nil {
		return &apimodel.MusicSet{}, err
	}

	if err := d.setTunesInApiSet(apiSet); err != nil {
		return &apimodel.MusicSet{}, err
	}

	return apiSet, nil
}

func apiSetFromDbSet(dbSet *model.MusicSet) (*apimodel.MusicSet, error) {
	apiSet := &apimodel.MusicSet{}
	if err := copier.Copy(apiSet, dbSet); err != nil {
		return nil, err
	}
	// Copier creates a default empty slice if the tunes are nil
	// so we set it to nil if there are no tunes
	if len(apiSet.Tunes) == 0 {
		apiSet.Tunes = nil
	}

	return apiSet, nil
}

func (d *dbService) getMusicSetByTuneIds(tuneIds []uuid.UUID) (*apimodel.MusicSet, error) {
	allMusicSets, err := d.MusicSets()
	if err != nil {
		return nil, err
	}

	var matchingSet *apimodel.MusicSet
	for _, set := range allMusicSets {
		if len(set.Tunes) != len(tuneIds) {
			continue
		}

		allTunesMatch := true
		for i, tune := range set.Tunes {
			if tuneIds[i] != tune.Id {
				allTunesMatch = false
				break
			}
		}
		if allTunesMatch {
			matchingSet = set
			break
		}
	}

	if matchingSet == nil {
		return nil, common.NotFound
	} else {
		return matchingSet, nil
	}
}

func (d *dbService) setTunesInApiSet(apiSet *apimodel.MusicSet) error {
	var setTunes []model.Tune
	err := d.db.Joins("JOIN music_set_tunes mst on tunes.id = mst.tune_id").
		Where("mst.music_set_id=?", apiSet.Id).
		Order("mst.\"order\"").
		Find(&setTunes).Error
	if err != nil {
		return fmt.Errorf("failed getting music set tunes: %s", err.Error())
	}

	if len(setTunes) > 0 {
		if err := copier.Copy(&apiSet.Tunes, &setTunes); err != nil {
			return err
		}
	}

	return nil
}

func (d *dbService) UpdateMusicSet(id uuid.UUID, updateSet apimodel.UpdateSet) (*apimodel.MusicSet, error) {
	if err := d.validator.ValidateUpdateSet(updateSet); err != nil {
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
	apiSet, err := apiSetFromDbSet(dbSet)
	if err != nil {
		return nil, err
	}

	return apiSet, nil
}

func (d *dbService) DeleteMusicSet(id uuid.UUID) error {
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
	setId uuid.UUID,
	tuneIds []uuid.UUID,
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
func (d *dbService) dbTunesFromIds(tuneIds []uuid.UUID) ([]model.Tune, error) {
	if len(tuneIds) == 0 {
		return nil, nil
	}

	var distinctTuneIds []uuid.UUID
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

func (d *dbService) assignMusicSetTunes(setId uuid.UUID, tuneIds []uuid.UUID) error {
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

func tunesOrderedByIds(tunes []model.Tune, tuneIds []uuid.UUID) []model.Tune {
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

func (d *dbService) ImportTunes(
	tunes []*messages.ImportedTune,
	fileInfo *common.ImportFileInfo,
) ([]*apimodel.ImportTune, *apimodel.BasicMusicSet, error) {
	var apiTunes []*apimodel.ImportTune
	var musicSet *apimodel.BasicMusicSet

	if fileInfo == nil {
		return nil, nil, fmt.Errorf("no file info given")
	}
	fileExists, err := d.hasImportFile(fileInfo)
	if err != nil {
		return nil, nil, err
	}

	if fileExists {
		return nil, nil, fmt.Errorf("file %s already imported", fileInfo.Name)
	}

	err = d.db.Transaction(func(tx *gorm.DB) error {
		importFile, err := d.createImportFile(fileInfo)
		if err != nil {
			return err
		}

		for _, impTune := range tunes {
			newTune := impTune.Tune
			timeSigStr := ""
			timeSig := newTune.FirstTimeSignature()
			if timeSig != nil {
				timeSigStr = timeSig.DisplayString()
			}

			if impTune.TuneFileData != nil {
				hasData, err := d.hasSingleFileData(impTune.TuneFileData)
				if err != nil {
					return err
				}

				if hasData {
					alreadyImportedTune, err := d.getTuneWithSingleFileData(impTune.TuneFileData)
					if err != nil {
						return err
					}

					if alreadyImportedTune != nil {
						impTune := &apimodel.ImportTune{}
						err = copier.Copy(impTune, alreadyImportedTune)
						if err != nil {
							return fmt.Errorf("failed creating import tune from already imported tune: %s", err.Error())
						}
						apiTunes = append(apiTunes, impTune)
						continue
					}
				}
			}

			createTune := apimodel.CreateTune{
				Title:    newTune.Title,
				Type:     newTune.Type,
				TimeSig:  timeSigStr,
				Composer: newTune.Composer,
				Arranger: newTune.Arranger,
			}
			apiTune, err := d.CreateTune(createTune, importFile)
			if err != nil {
				return err
			}
			muMoTuneFile, err := model.TuneFileFromMusicModelTune(newTune)
			if err != nil {
				return err
			}

			if err = d.AddFileToTune(apiTune.Id, muMoTuneFile); err != nil {
				return err
			}

			if impTune.TuneFileData != nil {
				muMoTuneFile = &model.TuneFile{
					Type:           fileInfo.FileType,
					Data:           impTune.TuneFileData,
					SingleTuneData: true,
				}
				if err = d.AddFileToTune(apiTune.Id, muMoTuneFile); err != nil {
					return err
				}
			}

			importTune := &apimodel.ImportTune{}
			err = copier.Copy(&importTune, apiTune)
			if err != nil {
				return err
			}
			setMessagesToApiTune(importTune, newTune)
			apiTunes = append(apiTunes, importTune)
		}

		if len(tunes) > 1 {
			var apiSet *apimodel.MusicSet
			var err error
			var tuneIds []uuid.UUID
			for _, apiTune := range apiTunes {
				tuneIds = append(tuneIds, apiTune.Id)
			}

			musicSetTitle := musicSetTitleFromTunes(apiTunes)
			apiSet, err = d.getMusicSetByTuneIds(tuneIds)
			if errors.Is(err, common.NotFound) { // music set not yet in db
				createSet := apimodel.CreateSet{
					Title: musicSetTitle,
					Tunes: tuneIds,
				}
				apiSet, err = d.CreateMusicSet(createSet, importFile)
				if err != nil {
					return fmt.Errorf("failed creating set for file %s: %s", fileInfo.Name, err.Error())
				}
			}

			musicSet = &apimodel.BasicMusicSet{}
			err = copier.Copy(musicSet, apiSet)
			if err != nil {
				return fmt.Errorf("failed creating basic music set from music set")
			}
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return apiTunes, musicSet, nil
}

func musicSetTitleFromTunes(tunes []*apimodel.ImportTune) string {
	if len(tunes) <= 1 {
		log.Error().Msg("not enough tunes given to create a music set title")
		return ""
	}

	if tuneTypesAndTimeSigsAreTheSame(tunes) {
		return fmt.Sprintf("%s %s Set", tunes[0].TimeSig, tunes[0].Type)
	}

	if tuneTypesAreTheSame(tunes) {
		return fmt.Sprintf("%s Set", tunes[0].Type)
	}

	if isMsr(tunes) {
		return "MSR Set"
	}

	var tuneTypes []string
	for _, tune := range tunes {
		if strings.TrimSpace(tune.Type) == "" {
			tuneTypes = append(tuneTypes, "Unknown Type")
		} else {
			tuneTypes = append(tuneTypes, tune.Type)
		}
	}

	return fmt.Sprint(strings.Join(tuneTypes, " - "))
}

func tuneTypesAndTimeSigsAreTheSame(tunes []*apimodel.ImportTune) bool {
	var timeSig, tuneType string
	for i, tune := range tunes {
		if i == 0 {
			timeSig = tune.TimeSig
			tuneType = tune.Type
			continue
		}

		if tune.TimeSig != timeSig ||
			tune.Type != tuneType {
			return false
		}
	}

	return true
}

func tuneTypesAreTheSame(tunes []*apimodel.ImportTune) bool {
	var tuneType string
	for i, tune := range tunes {
		if i == 0 {
			tuneType = tune.Type
			continue
		}

		if tune.Type != tuneType {
			return false
		}
	}

	return true
}

func isMsr(tunes []*apimodel.ImportTune) bool {
	if len(tunes) != 3 {
		return false
	}

	if strings.ToLower(tunes[0].Type) == "march" &&
		strings.ToLower(tunes[1].Type) == "strathspey" &&
		strings.ToLower(tunes[2].Type) == "reel" {
		return true
	}

	return false
}

// hasSingleFileData true if a tune with the same single tune file data exists in the database.
func (d *dbService) hasSingleFileData(
	data []byte,
) (bool, error) {
	if data == nil || len(data) == 0 {
		return false, nil
	}

	_, err := d.getSingleTuneFileByData(data)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// getTuneWithSingleFileData returns the tune with the same single tune file data in database
func (d *dbService) getTuneWithSingleFileData(
	data []byte,
) (*apimodel.Tune, error) {
	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("no file data given to search for")
	}

	tf, err := d.getSingleTuneFileByData(data)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // no tune with this data found => return nil for tune
	}
	if err != nil {
		return nil, err
	}

	ft, err := d.GetTune(tf.TuneID)
	if err != nil {
		return nil, err
	}

	return ft, nil
}

func (d *dbService) getSingleTuneFileByData(data []byte) (*model.TuneFile, error) {
	tf := &model.TuneFile{}
	if err := d.db.Where("data = ? AND single_tune_data IS TRUE", data).
		First(tf).Error; err != nil {
		return nil, err
	}

	return tf, nil
}

func setMessagesToApiTune(apiTune *apimodel.ImportTune, tune *tune.Tune) {
	im := tune.ImportMessages()
	for _, message := range im {
		switch message.Severity {
		case measure.Severity_Error:
			apiTune.Errors = append(apiTune.Errors, message.Text)
		case measure.Severity_Warning:
			apiTune.Warnings = append(apiTune.Warnings, message.Text)
		case measure.Severity_Info:
			apiTune.Infos = append(apiTune.Infos, message.Text)
		}
	}
}

func (d *dbService) GetTuneFile(tuneId uuid.UUID, fType file_type.Type) (*model.TuneFile, error) {
	tuneFile := &model.TuneFile{
		TuneID: tuneId,
		Type:   fType,
	}

	if err := d.db.First(tuneFile).Error; err != nil {
		return nil, common.NotFound
	}

	return tuneFile, nil
}

func (d *dbService) GetTuneFiles(tuneId uuid.UUID) ([]*model.TuneFile, error) {
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

func (d *dbService) AddFileToTune(tuneId uuid.UUID, tFile *model.TuneFile) error {
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

func (d *dbService) DeleteFileFromTune(tuneId uuid.UUID, fType file_type.Type) error {
	tuneFile := &model.TuneFile{
		TuneID: tuneId,
		Type:   fType,
	}
	if err := d.db.Delete(tuneFile).Error; err != nil {
		return err
	}

	return nil
}

func NewDbDataService(
	db *gorm.DB,
	validator interfaces.ApiModelValidator,
) interfaces.DataService {
	return &dbService{
		db:        db,
		validator: validator,
	}
}
