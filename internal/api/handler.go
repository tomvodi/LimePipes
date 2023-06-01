package api

import (
	"banduslib/internal/api/apimodel"
	"banduslib/internal/common"
	"banduslib/internal/common/music_model"
	"banduslib/internal/interfaces"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type apiHandler struct {
	service             interfaces.DataService
	bwwParser           interfaces.BwwParser
	bwwFileTuneSplitter interfaces.BwwFileByTuneSplitter
}

func (a *apiHandler) ImportBww(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]

	var importFiles []*apimodel.ImportFile
	for _, file := range files {
		importFile, err := a.importBwwFile(file)
		if err != nil {
			httpErrorResponse(c, http.StatusInternalServerError, err)
			return
		}

		importFiles = append(importFiles, importFile)
	}
	c.JSON(http.StatusOK, importFiles)
}

func (a *apiHandler) importBwwFile(
	file *multipart.FileHeader,
) (*apimodel.ImportFile, error) {
	fileReader, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed open file %s for reading", file.Filename)
	}
	defer fileReader.Close()

	fileData, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, fmt.Errorf("failed reading file %s: %s", file.Filename, err.Error())
	}

	importFile := &apimodel.ImportFile{
		Name: file.Filename,
		Result: apimodel.ParseResult{
			Success: false,
		},
	}

	bwwFileTuneData, err := a.bwwFileTuneSplitter.SplitFileData(fileData)
	if err != nil {
		return nil, err
	}

	var muModel music_model.MusicModel
	muModel, err = a.bwwParser.ParseBwwData(fileData)
	if err != nil {
		importFile.Result = apimodel.ParseResult{
			Message: err.Error(),
		}
		return importFile, err
	}

	fileName := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))

	apiImpTunes, err := a.service.ImportMusicModel(muModel, fileName, bwwFileTuneData)
	if err != nil {
		importFile.Result = apimodel.ParseResult{
			Message: err.Error(),
		}
		return importFile, err
	}

	importFile.Result.Success = true
	importFile.Tunes = apiImpTunes
	return importFile, nil
}

func httpErrorResponse(c *gin.Context, code int, err error) {
	c.JSON(code, apimodel.Error{
		Code:    int32(code),
		Message: err.Error(),
	})
}

func handleResponseForError(c *gin.Context, err error) {
	code := http.StatusInternalServerError
	if err == common.NotFound {
		code = http.StatusNotFound
	}

	c.JSON(code, apimodel.Error{
		Code:    int32(code),
		Message: err.Error(),
	})
}

func (a *apiHandler) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (a *apiHandler) CreateTune(c *gin.Context) {
	var createTune apimodel.CreateTune
	if err := c.ShouldBindJSON(&createTune); err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	tune, err := a.service.CreateTune(createTune)
	if err != nil {
		httpErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, tune)
}

func (a *apiHandler) GetTune(c *gin.Context) {
	tuneId, err := strconv.ParseUint(c.Param("tuneId"), 10, 64)
	if err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	tune, err := a.service.GetTune(tuneId)
	if err != nil {
		handleResponseForError(c, err)
		return
	}

	c.JSON(http.StatusOK, tune)
}

func (a *apiHandler) ListTunes(c *gin.Context) {
	tunes, err := a.service.Tunes()
	if err != nil {
		httpErrorResponse(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, tunes)
}

func (a *apiHandler) UpdateTune(c *gin.Context) {
	var updateTune apimodel.UpdateTune
	if err := c.ShouldBindJSON(&updateTune); err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	tuneId, err := strconv.ParseUint(c.Param("tuneId"), 10, 64)
	if err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	tune, err := a.service.UpdateTune(tuneId, updateTune)
	if err != nil {
		handleResponseForError(c, err)
		return
	}

	c.JSON(http.StatusOK, tune)
}

func (a *apiHandler) DeleteTune(c *gin.Context) {
	tuneId, err := strconv.ParseUint(c.Param("tuneId"), 10, 64)
	if err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := a.service.DeleteTune(tuneId); err != nil {
		handleResponseForError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *apiHandler) CreateSet(c *gin.Context) {
	var createSet apimodel.CreateSet

	if err := c.ShouldBindJSON(&createSet); err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	set, err := a.service.CreateMusicSet(createSet)
	if err != nil {
		httpErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, set)
}

func (a *apiHandler) GetSet(c *gin.Context) {
	setId, err := strconv.ParseUint(c.Param("setId"), 10, 64)
	if err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	set, err := a.service.GetMusicSet(setId)
	if err != nil {
		handleResponseForError(c, err)
		return
	}

	c.JSON(http.StatusOK, set)
}

func (a *apiHandler) ListSets(c *gin.Context) {
	sets, err := a.service.MusicSets()
	if err != nil {
		httpErrorResponse(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, sets)
}

func (a *apiHandler) UpdateSet(c *gin.Context) {
	var updateSet apimodel.UpdateSet
	if err := c.ShouldBindJSON(&updateSet); err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	setId, err := strconv.ParseUint(c.Param("setId"), 10, 64)
	if err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	set, err := a.service.UpdateMusicSet(setId, updateSet)
	if err != nil {
		handleResponseForError(c, err)
		return
	}

	c.JSON(http.StatusOK, set)
}

func (a *apiHandler) DeleteSet(c *gin.Context) {
	setId, err := strconv.ParseUint(c.Param("setId"), 10, 64)
	if err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := a.service.DeleteMusicSet(setId); err != nil {
		handleResponseForError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *apiHandler) AssignTunesToSet(c *gin.Context) {
	var tuneIds []uint64
	if err := c.ShouldBindJSON(&tuneIds); err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	setId, err := strconv.ParseUint(c.Param("setId"), 10, 64)
	if err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	set, err := a.service.AssignTunesToMusicSet(setId, tuneIds)
	if err != nil {
		handleResponseForError(c, err)
		return
	}

	c.JSON(http.StatusOK, set)
}

func NewApiHandler(
	service interfaces.DataService,
	bwwParser interfaces.BwwParser,
	bwwFileTuneSplitter interfaces.BwwFileByTuneSplitter,
) interfaces.ApiHandler {
	return &apiHandler{
		service:             service,
		bwwParser:           bwwParser,
		bwwFileTuneSplitter: bwwFileTuneSplitter,
	}
}
