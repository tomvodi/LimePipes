package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
	"github.com/tomvodi/limepipes/internal/apigen/apimodel"
	"github.com/tomvodi/limepipes/internal/common"
	"github.com/tomvodi/limepipes/internal/interfaces"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type Handler struct {
	service       interfaces.DataService
	pluginLoader  interfaces.PluginLoader
	healthChecker interfaces.HealthChecker
}

func (a *Handler) Home(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (a *Handler) Health(c *gin.Context) {
	handler, err := a.healthChecker.GetCheckHandler()
	if err != nil {
		httpErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	handleFunc := gin.WrapH(handler)
	handleFunc(c)
}

func (a *Handler) ImportFile(c *gin.Context) {
	iFile, err := c.FormFile("file")
	if err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	fExt := filepath.Ext(iFile.Filename)
	if fExt == "" {
		c.JSON(http.StatusBadRequest,
			apimodel.Error{
				Message: "import file does not have an extension",
			},
		)
		return
	}

	fType, err := a.pluginLoader.FileFormatForFileExtension(fExt)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			apimodel.Error{
				Message: fmt.Sprintf("file extension %s is currently not supported: %s", fExt, err.Error()),
			},
		)
		return
	}

	fInfo, err := a.createImportFileInfo(iFile, fType)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			apimodel.Error{
				Message: err.Error(),
			})
		return
	}
	_, err = a.service.GetImportFileByHash(fInfo.Hash)

	if !errors.Is(err, common.ErrNotFound) {
		c.JSON(http.StatusConflict,
			fmt.Sprintf("file %s was already imported", iFile.Filename))
		return
	}

	importTunes, importSet, err := a.importFile(fInfo, fExt)
	if err != nil {
		httpErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	importResponse := &apimodel.ImportFile{
		Name:  iFile.Filename,
		Set:   *importSet,
		Tunes: importTunes,
	}

	c.JSON(http.StatusOK, importResponse)
}

func (a *Handler) createImportFileInfo(
	iFile *multipart.FileHeader,
	fFormat fileformat.Format,
) (*common.ImportFileInfo, error) {
	fileReader, err := iFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed open file %s for reading", iFile.Filename)
	}
	defer fileReader.Close()

	fileData, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, fmt.Errorf("failed reading file %s: %s", iFile.Filename, err.Error())
	}

	fInfo, err := common.NewImportFileInfo(iFile.Filename, fFormat, fileData)
	if err != nil {
		return nil, fmt.Errorf("failed creating import file info for file %s: %s", iFile.Filename, err.Error())
	}

	return fInfo, nil
}

func (a *Handler) importFile(
	fInfo *common.ImportFileInfo,
	fileExt string,
) ([]*apimodel.ImportTune, *apimodel.BasicMusicSet, error) {
	filePlugin, err := a.pluginLoader.PluginForFileExtension(fileExt)
	if err != nil {
		return nil, nil, fmt.Errorf("fileData extension %s is currently not supported (no plugin): %s", fileExt, err.Error())
	}

	parsedTunes, err := filePlugin.Import(fInfo.Data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed parsing fileData %s: %s", fInfo.Name, err.Error())
	}

	return a.service.ImportTunes(parsedTunes.ImportedTunes, fInfo)
}

func httpErrorResponse(c *gin.Context, code int, err error) {
	c.JSON(code, apimodel.Error{
		Message: err.Error(),
	})
}

func handleResponseForError(c *gin.Context, err error) {
	code := http.StatusInternalServerError
	if errors.Is(err, common.ErrNotFound) {
		code = http.StatusNotFound
	}

	c.JSON(code, apimodel.Error{
		Message: err.Error(),
	})
}

func (a *Handler) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (a *Handler) CreateTune(c *gin.Context) {
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

func (a *Handler) GetTune(c *gin.Context) {
	tuneID, err := uuid.Parse(c.Param("tuneID"))
	if err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	tune, err := a.service.GetTune(tuneID)
	if err != nil {
		handleResponseForError(c, err)
		return
	}

	c.JSON(http.StatusOK, tune)
}

func (a *Handler) ListTunes(c *gin.Context) {
	tunes, err := a.service.Tunes()
	if err != nil {
		httpErrorResponse(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, tunes)
}

func (a *Handler) UpdateTune(c *gin.Context) {
	var updateTune apimodel.UpdateTune
	if err := c.ShouldBindJSON(&updateTune); err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	tuneID, err := uuid.Parse(c.Param("tuneID"))
	if err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	tune, err := a.service.UpdateTune(tuneID, updateTune)
	if err != nil {
		handleResponseForError(c, err)
		return
	}

	c.JSON(http.StatusOK, tune)
}

func (a *Handler) DeleteTune(c *gin.Context) {
	tuneID, err := uuid.Parse(c.Param("tuneID"))
	if err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := a.service.DeleteTune(tuneID); err != nil {
		handleResponseForError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *Handler) CreateSet(c *gin.Context) {
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

func (a *Handler) GetSet(c *gin.Context) {
	setID, err := uuid.Parse(c.Param("setID"))
	if err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	set, err := a.service.GetMusicSet(setID)
	if err != nil {
		handleResponseForError(c, err)
		return
	}

	c.JSON(http.StatusOK, set)
}

func (a *Handler) ListSets(c *gin.Context) {
	sets, err := a.service.MusicSets()
	if err != nil {
		httpErrorResponse(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, sets)
}

func (a *Handler) UpdateSet(c *gin.Context) {
	var updateSet apimodel.UpdateSet
	if err := c.ShouldBindJSON(&updateSet); err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	setID, err := uuid.Parse(c.Param("setID"))
	if err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	set, err := a.service.UpdateMusicSet(setID, updateSet)
	if err != nil {
		handleResponseForError(c, err)
		return
	}

	c.JSON(http.StatusOK, set)
}

func (a *Handler) DeleteSet(c *gin.Context) {
	setID, err := uuid.Parse(c.Param("setID"))
	if err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := a.service.DeleteMusicSet(setID); err != nil {
		handleResponseForError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *Handler) AssignTunesToSet(c *gin.Context) {
	var tuneIDs []uuid.UUID
	if err := c.ShouldBindJSON(&tuneIDs); err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	setID, err := uuid.Parse(c.Param("setID"))
	if err != nil {
		httpErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	set, err := a.service.AssignTunesToMusicSet(setID, tuneIDs)
	if err != nil {
		handleResponseForError(c, err)
		return
	}

	c.JSON(http.StatusOK, set)
}

func NewAPIHandler(
	service interfaces.DataService,
	pluginLoader interfaces.PluginLoader,
	healthChecker interfaces.HealthChecker,
) *Handler {
	return &Handler{
		service:       service,
		pluginLoader:  pluginLoader,
		healthChecker: healthChecker,
	}
}
