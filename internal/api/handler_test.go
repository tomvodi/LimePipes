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
	"github.com/tomvodi/limepipes/internal/api_gen/apimodel"
	"github.com/tomvodi/limepipes/internal/database/model"
	"github.com/tomvodi/limepipes/internal/interfaces/mocks"
	"github.com/tomvodi/limepipes/internal/utils"
	"io"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Api handler", func() {
	utils.SetupConsoleLogger()
	var c *gin.Context
	var httpRec *httptest.ResponseRecorder
	var api *apiHandler
	var testId1 uuid.UUID
	var dataService *mocks.DataService

	BeforeEach(func() {
		testId1 = uuid.MustParse("00000000-0000-0000-0000-000000000001")

		httpRec = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(httpRec)
		dataService = mocks.NewDataService(GinkgoT())
		api = &apiHandler{
			service: dataService,
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
				MockJsonPost(c, http.MethodPost, tune)
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
							Id:    testId1,
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
		var tuneIds []string
		var tuneIdsUuid []uuid.UUID
		var setId uuid.UUID

		JustBeforeEach(func() {
			api.AssignTunesToSet(c)
		})

		BeforeEach(func() {
			setId = testId1
			tuneIds = []string{
				"00000000-0000-0000-0000-000000000002",
				"00000000-0000-0000-0000-000000000003",
			}
			tuneIdsUuid = []uuid.UUID{
				uuid.MustParse(tuneIds[0]),
				uuid.MustParse(tuneIds[1]),
			}
			MockJsonPost(c, http.MethodPut, tuneIds)

			c.Params = gin.Params{
				{Key: "setId", Value: setId.String()},
			}
		})

		When("service successfully assignes tunes", func() {
			BeforeEach(func() {
				dataService.EXPECT().AssignTunesToMusicSet(setId, tuneIdsUuid).
					Return(&apimodel.MusicSet{
						Id:    setId,
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
})

func MockJsonPost(c *gin.Context, method string, content interface{}) {
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
