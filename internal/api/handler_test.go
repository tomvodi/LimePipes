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
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
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

var _ = Describe("Api Handler", func() {
	utils.SetupConsoleLogger()
	var c *gin.Context
	var httpRec *httptest.ResponseRecorder
	var api *Handler
	var testID1 uuid.UUID
	var dataService *mocks.DataService
	var healthChecker *mocks.HealthChecker
	var pluginLoader *mocks.PluginLoader
	var lpPlugin *pmocks.LimePipesPlugin

	BeforeEach(func() {
		testID1 = uuid.MustParse("00000000-0000-0000-0000-000000000001")

		httpRec = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(httpRec)
		dataService = mocks.NewDataService(GinkgoT())
		healthChecker = mocks.NewHealthChecker(GinkgoT())
		pluginLoader = mocks.NewPluginLoader(GinkgoT())
		lpPlugin = pmocks.NewLimePipesPlugin(GinkgoT())
		api = &Handler{
			service:       dataService,
			healthChecker: healthChecker,
			pluginLoader:  pluginLoader,
		}
	})

	Context("Home", func() {
		JustBeforeEach(func() {
			api.Home(c)
		})

		It("should return ok", func() {
			Expect(httpRec.Code).To(Equal(http.StatusOK))
		})
	})

	Context("Health", func() {
		JustBeforeEach(func() {
			api.Health(c)
		})

		When("health checker returns an error", func() {
			BeforeEach(func() {
				healthChecker.EXPECT().GetCheckHandler().
					Return(nil, fmt.Errorf("xxx"))
			})

			It("should return a server error", func() {
				Expect(httpRec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("health checker returns a handler", func() {
			BeforeEach(func() {
				healthChecker.EXPECT().GetCheckHandler().
					Return(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.WriteHeader(http.StatusOK)
					}), nil)
			})

			It("should return ok", func() {
				Expect(httpRec.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Context("Get Tune", func() {
		var tuneID uuid.UUID

		JustBeforeEach(func() {
			api.GetTune(c)
		})

		BeforeEach(func() {
			tuneID = testID1
			c.Params = gin.Params{
				{Key: "tuneID", Value: tuneID.String()},
			}
		})

		When("no uuid as tuneID", func() {
			BeforeEach(func() {
				c.Params = gin.Params{
					{Key: "tuneID", Value: "not a uuid"},
				}
			})

			It("should return BadRequest", func() {
				Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("service returns an error", func() {
			BeforeEach(func() {
				dataService.EXPECT().GetTune(tuneID).
					Return(nil, fmt.Errorf("xxx"))
			})

			It("should return a server error", func() {
				Expect(httpRec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("service returns a tune", func() {
			BeforeEach(func() {
				dataService.EXPECT().GetTune(tuneID).
					Return(&apimodel.Tune{
						Id:    tuneID,
						Title: "test title",
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

	Context("List Tunes", func() {
		JustBeforeEach(func() {
			api.ListTunes(c)
		})

		When("service returns an error", func() {
			BeforeEach(func() {
				dataService.EXPECT().Tunes().
					Return(nil, fmt.Errorf("xxx"))
			})

			It("should return a server error", func() {
				Expect(httpRec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("service returns tunes", func() {
			BeforeEach(func() {
				dataService.EXPECT().Tunes().
					Return([]*apimodel.Tune{
						{
							Id:    testID1,
							Title: "test title",
						},
					}, nil)
			})

			It("should return ok and the tunes", func() {
				Expect(httpRec.Code).To(Equal(http.StatusOK))
				data, err := io.ReadAll(httpRec.Body)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(data)).To(Equal("[{\"id\":\"00000000-0000-0000-0000-000000000001\",\"title\":\"test title\"}]"))
			})
		})
	})

	Context("Update Tune", func() {
		var tuneID uuid.UUID

		JustBeforeEach(func() {
			api.UpdateTune(c)
		})

		BeforeEach(func() {
			tuneID = testID1
			c.Params = gin.Params{
				{Key: "tuneID", Value: tuneID.String()},
			}
		})

		When("updating a tune without a title", func() {
			var tune apimodel.UpdateTune

			BeforeEach(func() {
				tune = apimodel.UpdateTune{
					Title: "",
				}
				mockJSONPost(c, http.MethodPut, tune)
			})

			It("should return BadRequest", func() {
				Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("updating a tune", func() {
			var tune apimodel.UpdateTune

			BeforeEach(func() {
				tune = apimodel.UpdateTune{
					Title: "test title",
				}
				mockJSONPost(c, http.MethodPut, tune)
			})

			When("no uuid as tuneID", func() {
				BeforeEach(func() {
					c.Params = gin.Params{
						{Key: "tuneID", Value: "not a uuid"},
					}
				})

				It("should return BadRequest", func() {
					Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
				})
			})

			When("service returns an error on update", func() {
				BeforeEach(func() {
					dataService.EXPECT().UpdateTune(tuneID, tune).
						Return(nil, fmt.Errorf("xxx"))
				})

				It("should return a server error", func() {
					Expect(httpRec.Code).To(Equal(http.StatusInternalServerError))
				})
			})

			When("service successfully updates tune", func() {
				BeforeEach(func() {
					dataService.EXPECT().UpdateTune(tuneID, tune).
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

	Context("Delete Tune", func() {
		var tuneID uuid.UUID

		JustBeforeEach(func() {
			api.DeleteTune(c)
		})

		BeforeEach(func() {
			tuneID = testID1
			c.Params = gin.Params{
				{Key: "tuneID", Value: tuneID.String()},
			}
		})

		When("no uuid as tuneID", func() {
			BeforeEach(func() {
				c.Params = gin.Params{
					{Key: "tuneID", Value: "not a uuid"},
				}
			})

			It("should return BadRequest", func() {
				Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("service returns an error", func() {
			BeforeEach(func() {
				dataService.EXPECT().DeleteTune(tuneID).
					Return(fmt.Errorf("xxx"))
			})

			It("should return a server error", func() {
				Expect(httpRec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("service successfully deletes tune", func() {
			BeforeEach(func() {
				dataService.EXPECT().DeleteTune(tuneID).
					Return(nil)
			})

			It("should return ok", func() {
				Expect(httpRec.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Context("Create Tune", func() {
		JustBeforeEach(func() {
			api.CreateTune(c)
		})

		Context("no data given", func() {
			It("should return BadRequest", func() {
				Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("creating a tune without a title", func() {
			var tune apimodel.CreateTune

			BeforeEach(func() {
				tune = apimodel.CreateTune{
					Title: "",
				}
				mockJSONPost(c, http.MethodPost, tune)
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

	Context("Get Set", func() {
		var setID uuid.UUID

		JustBeforeEach(func() {
			api.GetSet(c)
		})

		BeforeEach(func() {
			setID = testID1
			c.Params = gin.Params{
				{Key: "setID", Value: setID.String()},
			}
		})

		When("no uuid as setID", func() {
			BeforeEach(func() {
				c.Params = gin.Params{
					{Key: "setID", Value: "not a uuid"},
				}
			})

			It("should return BadRequest", func() {
				Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("service returns an error", func() {
			BeforeEach(func() {
				dataService.EXPECT().GetMusicSet(setID).
					Return(nil, fmt.Errorf("xxx"))
			})

			It("should return a server error", func() {
				Expect(httpRec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("service returns a set", func() {
			BeforeEach(func() {
				dataService.EXPECT().GetMusicSet(setID).
					Return(&apimodel.MusicSet{
						Id:    setID,
						Title: "test title",
					}, nil)
			})

			It("should return ok and the set", func() {
				Expect(httpRec.Code).To(Equal(http.StatusOK))
				data, err := io.ReadAll(httpRec.Body)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(data)).To(Equal("{\"id\":\"00000000-0000-0000-0000-000000000001\",\"title\":\"test title\"}"))
			})
		})
	})

	Context("Create Set", func() {
		JustBeforeEach(func() {
			api.CreateSet(c)
		})

		Context("no data given", func() {
			It("should return BadRequest", func() {
				Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("creating a set without a title", func() {
			var set apimodel.CreateSet

			BeforeEach(func() {
				set = apimodel.CreateSet{
					Title: "",
				}
				mockJSONPost(c, http.MethodPost, set)
			})

			It("should return BadRequest", func() {
				Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("creating a set", func() {
			var set apimodel.CreateSet

			BeforeEach(func() {
				set = apimodel.CreateSet{
					Title: "test title",
				}
				mockJSONPost(c, http.MethodPost, set)
			})

			When("service returns an error on creation", func() {
				BeforeEach(func() {
					dataService.EXPECT().CreateMusicSet(set, (*model.ImportFile)(nil)).
						Return(nil, fmt.Errorf("xxx"))
				})

				It("should return a server error", func() {
					Expect(httpRec.Code).To(Equal(http.StatusInternalServerError))
				})
			})

			When("service successfully creates set", func() {
				BeforeEach(func() {
					dataService.EXPECT().CreateMusicSet(set, (*model.ImportFile)(nil)).
						Return(&apimodel.MusicSet{
							Id:    testID1,
							Title: set.Title,
						}, nil)
				})

				It("should return ok and the set", func() {
					Expect(httpRec.Code).To(Equal(http.StatusOK))
					data, err := io.ReadAll(httpRec.Body)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(string(data)).To(Equal("{\"id\":\"00000000-0000-0000-0000-000000000001\",\"title\":\"test title\"}"))
				})
			})
		})
	})

	Context("List Sets", func() {
		JustBeforeEach(func() {
			api.ListSets(c)
		})

		When("service returns an error", func() {
			BeforeEach(func() {
				dataService.EXPECT().MusicSets().
					Return(nil, fmt.Errorf("xxx"))
			})

			It("should return a server error", func() {
				Expect(httpRec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("service returns sets", func() {
			BeforeEach(func() {
				dataService.EXPECT().MusicSets().
					Return([]*apimodel.MusicSet{
						{
							Id:    testID1,
							Title: "test title",
						},
					}, nil)
			})

			It("should return ok and the sets", func() {
				Expect(httpRec.Code).To(Equal(http.StatusOK))
				data, err := io.ReadAll(httpRec.Body)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(data)).To(Equal("[{\"id\":\"00000000-0000-0000-0000-000000000001\",\"title\":\"test title\"}]"))
			})
		})
	})

	Context("Update Set", func() {
		var setID uuid.UUID

		JustBeforeEach(func() {
			api.UpdateSet(c)
		})

		BeforeEach(func() {
			setID = testID1
			c.Params = gin.Params{
				{Key: "setID", Value: setID.String()},
			}
		})

		When("updating a set without a title", func() {
			var set apimodel.UpdateSet

			BeforeEach(func() {
				set = apimodel.UpdateSet{
					Title: "",
				}
				mockJSONPost(c, http.MethodPut, set)
			})

			It("should return BadRequest", func() {
				Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("updating a set", func() {
			var set apimodel.UpdateSet

			BeforeEach(func() {
				set = apimodel.UpdateSet{
					Title: "test title",
				}
				mockJSONPost(c, http.MethodPut, set)
			})

			When("no uuid as setID", func() {
				BeforeEach(func() {
					c.Params = gin.Params{
						{Key: "setID", Value: "not a uuid"},
					}
				})

				It("should return BadRequest", func() {
					Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
				})
			})

			When("service returns an error on update", func() {
				BeforeEach(func() {
					dataService.EXPECT().UpdateMusicSet(setID, set).
						Return(nil, fmt.Errorf("xxx"))
				})

				It("should return a server error", func() {
					Expect(httpRec.Code).To(Equal(http.StatusInternalServerError))
				})
			})

			When("service successfully updates set", func() {
				BeforeEach(func() {
					dataService.EXPECT().UpdateMusicSet(setID, set).
						Return(&apimodel.MusicSet{
							Id:    testID1,
							Title: set.Title,
						}, nil)
				})

				It("should return ok and the set", func() {
					Expect(httpRec.Code).To(Equal(http.StatusOK))
					data, err := io.ReadAll(httpRec.Body)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(string(data)).To(Equal("{\"id\":\"00000000-0000-0000-0000-000000000001\",\"title\":\"test title\"}"))
				})
			})
		})
	})

	Context("Delete Set", func() {
		var setID uuid.UUID

		JustBeforeEach(func() {
			api.DeleteSet(c)
		})

		BeforeEach(func() {
			setID = testID1
			c.Params = gin.Params{
				{Key: "setID", Value: setID.String()},
			}
		})

		When("no uuid as setID", func() {
			BeforeEach(func() {
				c.Params = gin.Params{
					{Key: "setID", Value: "not a uuid"},
				}
			})

			It("should return BadRequest", func() {
				Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("service returns an error", func() {
			BeforeEach(func() {
				dataService.EXPECT().DeleteMusicSet(setID).
					Return(fmt.Errorf("xxx"))
			})

			It("should return a server error", func() {
				Expect(httpRec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("service successfully deletes set", func() {
			BeforeEach(func() {
				dataService.EXPECT().DeleteMusicSet(setID).
					Return(nil)
			})

			It("should return ok", func() {
				Expect(httpRec.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Context("Assign Tunes To Set", func() {
		var setID uuid.UUID
		var testID2 uuid.UUID

		JustBeforeEach(func() {
			api.AssignTunesToSet(c)
		})

		BeforeEach(func() {
			setID = testID1
			testID2 = uuid.MustParse("00000000-0000-0000-0000-000000000002")
			c.Params = gin.Params{
				{Key: "setID", Value: setID.String()},
			}
		})

		When("no valid tune IDs given", func() {
			BeforeEach(func() {
				mockJSONPost(c, http.MethodPost, []string{"xxx", "yyy"})
			})

			It("should return BadRequest", func() {
				Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("valid tune IDs given", func() {
			BeforeEach(func() {
				mockJSONPost(c, http.MethodPost, []uuid.UUID{
					testID1,
					testID2,
				})
			})

			When("no uuid as setID", func() {
				BeforeEach(func() {
					c.Params = gin.Params{
						{Key: "setID", Value: "not a uuid"},
					}
				})

				It("should return BadRequest", func() {
					Expect(httpRec.Code).To(Equal(http.StatusBadRequest))
				})
			})

			When("service returns an error", func() {
				BeforeEach(func() {
					dataService.EXPECT().AssignTunesToMusicSet(setID, []uuid.UUID{testID1, testID2}).
						Return(nil, fmt.Errorf("xxx"))
				})

				It("should return a server error", func() {
					Expect(httpRec.Code).To(Equal(http.StatusInternalServerError))
				})
			})

			When("service successfully assigns tunes to set", func() {
				BeforeEach(func() {
					dataService.EXPECT().AssignTunesToMusicSet(setID, []uuid.UUID{testID1, testID2}).
						Return(&apimodel.MusicSet{
							Id:    testID1,
							Title: "test music set",
						}, nil)
				})

				It("should return ok and the set", func() {
					Expect(httpRec.Code).To(Equal(http.StatusOK))
					data, err := io.ReadAll(httpRec.Body)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(string(data)).To(Equal("{\"id\":\"00000000-0000-0000-0000-000000000001\",\"title\":\"test music set\"}"))
				})
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
				pluginLoader.EXPECT().FileFormatForFileExtension(".abc").
					Return(fileformat.Format_Unknown, fmt.Errorf("file extension .abc is not supported"))
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
				pluginLoader.EXPECT().FileFormatForFileExtension(".bww").
					Return(fileformat.Format_BWW, nil)
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
				pluginLoader.EXPECT().FileFormatForFileExtension(".bww").
					Return(fileformat.Format_BWW, nil)
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
				pluginLoader.EXPECT().FileFormatForFileExtension(".bww").
					Return(fileformat.Format_BWW, nil)
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
				pluginLoader.EXPECT().FileFormatForFileExtension(".bww").
					Return(fileformat.Format_BWW, nil)
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
