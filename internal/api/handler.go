package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes/internal/api_gen/apimodel"
	api_interfaces "github.com/tomvodi/limepipes/internal/api_gen/interfaces"
	"github.com/tomvodi/limepipes/internal/common"
	"github.com/tomvodi/limepipes/internal/interfaces"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
)

type apiHandler struct {
	service       interfaces.DataService
	pluginLoader  interfaces.PluginLoader
	healthChecker interfaces.HealthChecker
}

func (a *apiHandler) Home(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (a *apiHandler) Health(c *gin.Context) {
	handler, err := a.healthChecker.GetCheckHandler()
	if err != nil {
		httpErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleFunc := gin.WrapH(handler)
	handleFunc(c)
}

func (a *apiHandler) ImportBww(c *gin.Context) {
	logReq, err := httputil.DumpRequest(c.Request, true)
	if err != nil {
		log.Err(err).Msg("failed dumping request")
	}
	log.Info().Msgf("request: %s", string(logReq))

	form, err := c.MultipartForm()
	if err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	var allFiles []*multipart.FileHeader
	for _, fh := range form.File {
		allFiles = append(allFiles, fh...)
	}

	var importFiles []*apimodel.ImportFile
	for _, file := range allFiles {
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

	bwwPlugin, err := a.pluginLoader.PluginForFileExtension(".bww")
	if err != nil {
		importFile.Result = apimodel.ParseResult{
			Message: err.Error(),
		}
		return importFile, err
	}

	importTunes, err := bwwPlugin.Import(fileData)
	if err != nil {
		importFile.Result = apimodel.ParseResult{
			Message: err.Error(),
		}
		return importFile, err
	}

	info, err := common.NewImportFileInfo(file.Filename, fileData)
	if err != nil {
		return nil, err
	}

	apiImpTunes, err := a.service.ImportTunes(importTunes.ImportedTunes, info)
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
	if errors.Is(err, common.NotFound) {
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

	tune, err := a.service.CreateTune(createTune, nil)
	if err != nil {
		httpErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, tune)
}

func (a *apiHandler) GetTune(c *gin.Context) {
	tuneId, err := uuid.Parse(c.Param("tuneId"))
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

	tuneId, err := uuid.Parse(c.Param("tuneId"))
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
	tuneId, err := uuid.Parse(c.Param("tuneId"))
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

	set, err := a.service.CreateMusicSet(createSet, nil)
	if err != nil {
		httpErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, set)
}

func (a *apiHandler) GetSet(c *gin.Context) {
	setId, err := uuid.Parse(c.Param("setId"))
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

	setId, err := uuid.Parse(c.Param("setId"))
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
	setId, err := uuid.Parse(c.Param("setId"))
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
	var tuneIds []uuid.UUID
	if err := c.ShouldBindJSON(&tuneIds); err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	setId, err := uuid.Parse(c.Param("setId"))
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
	pluginLoader interfaces.PluginLoader,
	healthChecker interfaces.HealthChecker,
) api_interfaces.ApiHandler {
	return &apiHandler{
		service:       service,
		pluginLoader:  pluginLoader,
		healthChecker: healthChecker,
	}
}
