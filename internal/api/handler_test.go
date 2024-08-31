package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/tune"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/file_type"
	pmocks "github.com/tomvodi/limepipes-plugin-api/plugin/v1/interfaces/mocks"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"github.com/tomvodi/limepipes/internal/apigen/apimodel"
	"github.com/tomvodi/limepipes/internal/common"
	"github.com/tomvodi/limepipes/internal/database/model"
	"github.com/tomvodi/limepipes/internal/interfaces/mocks"
	"github.com/tomvodi/limepipes/internal/utils"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
)

type multipartRequest struct {
	Fieldname  string
	Filename   string
	Content    []byte
	Endpoint   string
	HTTPMethod string
}

var _ = Describe("Api handler", func() {
	utils.SetupConsoleLogger()
	var c *gin.Context
	var httpRec *httptest.ResponseRecorder
	var api *apiHandler
	var testID1 uuid.UUID
	var dataService *mocks.DataService
	var pluginLoader *mocks.PluginLoader
	var lpPlugin *pmocks.LimePipesPlugin

	BeforeEach(func() {
		testID1 = uuid.MustParse("00000000-0000-0000-0000-000000000001")

		httpRec = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(httpRec)
		dataService = mocks.NewDataService(GinkgoT())
		pluginLoader = mocks.NewPluginLoader(GinkgoT())
		lpPlugin = pmocks.NewLimePipesPlugin(GinkgoT())
		api = &apiHandler{
			service:      dataService,
			pluginLoader: pluginLoader,
		}
	})

	Context("Create Tune", func() {
		JustBeforeEach(func() {
			api.CreateTune(c)
		})

		Context("no data given", func() {
			BeforeEach(func() {
				api.CreateTune(c)
			})

			It("should return BadRequest", func() {
				Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("creating a tune", func() {
			var tune apimodel.CreateTune

			BeforeEach(func() {
				tune = apimodel.CreateTune{
					Title: "test title",
				}
				mockJSONPost(c, http.MethodPost, tune)
			})

			When("service returns an error on creation", func() {
				BeforeEach(func() {
					dataService.EXPECT().CreateTune(tune, (*model.ImportFile)(nil)).
						Return(nil, fmt.Errorf("xxx"))
				})

				It("should return a server error", func() {
					Expect(httpRec.Code).To(Equal(http.StatusInternalServerError))
				})
			})

			When("service successfully creates tune", func() {
				BeforeEach(func() {
					dataService.EXPECT().CreateTune(tune, (*model.ImportFile)(nil)).
						Return(&apimodel.Tune{
							Id:    testID1,
							Title: tune.Title,
						}, nil)
				})

				It("should return ok and the tune", func() {
					Expect(httpRec.Code).To(Equal(http.StatusOK))
					data, err := io.ReadAll(httpRec.Body)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(string(data)).To(Equal("{\"id\":\"00000000-0000-0000-0000-000000000001\",\"title\":\"test title\"}"))
				})
			})
		})
	})

	Context("Assign tunes to a set", func() {
		var tuneIDs []string
		var tuneIDsUUID []uuid.UUID
		var setID uuid.UUID

		JustBeforeEach(func() {
			api.AssignTunesToSet(c)
		})

		BeforeEach(func() {
			setID = testID1
			tuneIDs = []string{
				"00000000-0000-0000-0000-000000000002",
				"00000000-0000-0000-0000-000000000003",
			}
			tuneIDsUUID = []uuid.UUID{
				uuid.MustParse(tuneIDs[0]),
				uuid.MustParse(tuneIDs[1]),
			}
			mockJSONPost(c, http.MethodPut, tuneIDs)

			c.Params = gin.Params{
				{Key: "setID", Value: setID.String()},
			}
		})

		When("service successfully assignes tunes", func() {
			BeforeEach(func() {
				dataService.EXPECT().AssignTunesToMusicSet(setID, tuneIDsUUID).
					Return(&apimodel.MusicSet{
						Id:    setID,
						Title: "set 1",
					}, nil)
			})

			It("should return ok and the tune", func() {
				Expect(httpRec.Code).To(Equal(http.StatusOK))
				data, err := io.ReadAll(httpRec.Body)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(data)).To(Equal("{\"id\":\"00000000-0000-0000-0000-000000000001\",\"title\":\"set 1\"}"))
			})
		})
	})

	Context("ImportFile", func() {
		JustBeforeEach(func() {
			api.ImportFile(c)
		})

		When("fieldname is wrong", func() {
			BeforeEach(func() {
				c.Request = multipartRequestForFile(multipartRequest{
					Fieldname:  "wrongfieldname",
					Filename:   "test.bww",
					Content:    []byte("test file content"),
					Endpoint:   "/imports",
					HTTPMethod: http.MethodPost,
				})
			})

			It("should return BadRequest", func() {
				api.ImportFile(c)
				Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("file has no extension", func() {
			BeforeEach(func() {
				c.Request = multipartRequestForFile(multipartRequest{
					Fieldname:  "file",
					Filename:   "test",
					Content:    []byte("test file content"),
					Endpoint:   "/imports",
					HTTPMethod: http.MethodPost,
				})
			})

			It("should return BadRequest", func() {
				api.ImportFile(c)
				Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("file extension is not supported", func() {
			BeforeEach(func() {
				c.Request = multipartRequestForFile(multipartRequest{
					Fieldname:  "file",
					Filename:   "test.abc",
					Content:    []byte("test file content"),
					Endpoint:   "/imports",
					HTTPMethod: http.MethodPost,
				})
				pluginLoader.EXPECT().FileTypeForFileExtension(".abc").
					Return(file_type.Type_Unknown, fmt.Errorf("file extension .abc is not supported"))
			})

			It("should return BadRequest", func() {
				api.ImportFile(c)
				Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("file was already imported", func() {
			BeforeEach(func() {
				c.Request = multipartRequestForFile(multipartRequest{
					Fieldname:  "file",
					Filename:   "test.bww",
					Content:    []byte("test file content"),
					Endpoint:   "/imports",
					HTTPMethod: http.MethodPost,
				})
				pluginLoader.EXPECT().FileTypeForFileExtension(".bww").
					Return(file_type.Type_BWW, nil)
				dataService.EXPECT().GetImportFileByHash("60f5237ed4049f0382661ef009d2bc42e48c3ceb3edb6600f7024e7ab3b838f3").
					Return(nil, nil)
			})

			It("should return a http conflict", func() {
				api.ImportFile(c)
				Expect(httpRec.Code).To(Equal(http.StatusConflict))
			})
		})

		When("there is no plugin for the file extension", func() {
			BeforeEach(func() {
				c.Request = multipartRequestForFile(multipartRequest{
					Fieldname:  "file",
					Filename:   "test.bww",
					Content:    []byte("test file content"),
					Endpoint:   "/imports",
					HTTPMethod: http.MethodPost,
				})
				pluginLoader.EXPECT().FileTypeForFileExtension(".bww").
					Return(file_type.Type_BWW, nil)
				dataService.EXPECT().GetImportFileByHash("60f5237ed4049f0382661ef009d2bc42e48c3ceb3edb6600f7024e7ab3b838f3").
					Return(nil, common.ErrNotFound)
				pluginLoader.EXPECT().PluginForFileExtension(".bww").
					Return(nil, fmt.Errorf("no plugin found for file extension .bww"))
			})

			It("should return InternalServerError", func() {
				api.ImportFile(c)
				Expect(httpRec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("plugin fails parsing the file", func() {
			BeforeEach(func() {
				c.Request = multipartRequestForFile(multipartRequest{
					Fieldname:  "file",
					Filename:   "test.bww",
					Content:    []byte("test file content"),
					Endpoint:   "/imports",
					HTTPMethod: http.MethodPost,
				})
				pluginLoader.EXPECT().FileTypeForFileExtension(".bww").
					Return(file_type.Type_BWW, nil)
				dataService.EXPECT().GetImportFileByHash("60f5237ed4049f0382661ef009d2bc42e48c3ceb3edb6600f7024e7ab3b838f3").
					Return(nil, common.ErrNotFound)
				pluginLoader.EXPECT().PluginForFileExtension(".bww").
					Return(lpPlugin, nil)
				lpPlugin.EXPECT().Import([]byte("test file content")).
					Return(nil, fmt.Errorf("failed parsing file test.bww: xxx"))
			})

			It("should return InternalServerError", func() {
				api.ImportFile(c)
				Expect(httpRec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("plugin successfully parses the file", func() {
			BeforeEach(func() {
				c.Request = multipartRequestForFile(multipartRequest{
					Fieldname:  "file",
					Filename:   "test.bww",
					Content:    []byte("test file content"),
					Endpoint:   "/imports",
					HTTPMethod: http.MethodPost,
				})
				pluginLoader.EXPECT().FileTypeForFileExtension(".bww").
					Return(file_type.Type_BWW, nil)
				dataService.EXPECT().GetImportFileByHash("60f5237ed4049f0382661ef009d2bc42e48c3ceb3edb6600f7024e7ab3b838f3").
					Return(nil, common.ErrNotFound)
				pluginLoader.EXPECT().PluginForFileExtension(".bww").
					Return(lpPlugin, nil)
				importTunes := []*messages.ImportedTune{
					{
						Tune: &tune.Tune{
							Title: "test title",
						},
						TuneFileData: []byte("test file content"),
					},
				}
				lpPlugin.EXPECT().Import([]byte("test file content")).
					Return(&messages.ImportFileResponse{
						ImportedTunes: importTunes,
					}, nil)
				dataService.EXPECT().ImportTunes(importTunes, mock.Anything).
					Return([]*apimodel.ImportTune{
						{
							Id:    testID1,
							Title: "test tune",
						},
					},
						&apimodel.BasicMusicSet{
							Id:    testID1,
							Title: "test music set",
						}, nil)
			})

			It("should return ok and the imported tunes", func() {
				api.ImportFile(c)
				Expect(httpRec.Code).To(Equal(http.StatusOK))
				data, err := io.ReadAll(httpRec.Body)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(data)).To(Equal("{\"name\":\"test.bww\",\"set\":{\"id\":\"00000000-0000-0000-0000-000000000001\",\"title\":\"test music set\"},\"tunes\":[{\"id\":\"00000000-0000-0000-0000-000000000001\",\"title\":\"test tune\"}]}{\"name\":\"test.bww\",\"set\":{\"id\":\"00000000-0000-0000-0000-000000000001\",\"title\":\"test music set\"},\"tunes\":[{\"id\":\"00000000-0000-0000-0000-000000000001\",\"title\":\"test tune\"}]}"))
			})
		})
	})
})

func multipartRequestForFile(
	req multipartRequest,
) *http.Request {
	body := new(bytes.Buffer)
	mulWriter := multipart.NewWriter(body)
	dataPart, err := mulWriter.CreateFormFile(req.Fieldname, req.Filename)
	Expect(err).NotTo(HaveOccurred())

	_, err = dataPart.Write(req.Content)
	Expect(err).NotTo(HaveOccurred())
	err = mulWriter.Close()
	Expect(err).NotTo(HaveOccurred())

	r, err := http.NewRequest(req.HTTPMethod, req.Endpoint, body)
	Expect(err).NotTo(HaveOccurred())
	r.Header.Set("Content-Type", mulWriter.FormDataContentType())

	return r
}

func mockJSONPost(c *gin.Context, method string, content any) {
	c.Request = &http.Request{
		Method: method,
		Header: make(http.Header),
	}
	c.Request.Header.Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(content)
	if err != nil {
		log.Error().Err(err).Msgf("failed marshal data %v", content)
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))
}
