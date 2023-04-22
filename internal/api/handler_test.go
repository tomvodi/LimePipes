package api

import (
	"banduslib/internal/api/apimodel"
	mock_interfaces "banduslib/internal/interfaces/mocks"
	"banduslib/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("CreateTune", func() {
	utils.SetupConsoleLogger()
	var c *gin.Context
	var httpRec *httptest.ResponseRecorder
	var api *apiHandler
	var mockCtrl *gomock.Controller
	var tune apimodel.CreateTune

	var dataService *mock_interfaces.MockDataService

	JustBeforeEach(func() {
		api.CreateTune(c)
	})

	BeforeEach(func() {
		httpRec = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(httpRec)
		mockCtrl = gomock.NewController(GinkgoT())
		dataService = mock_interfaces.NewMockDataService(mockCtrl)
		api = &apiHandler{
			service: dataService,
		}
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
		BeforeEach(func() {
			tune = apimodel.CreateTune{
				Title: "test title",
			}
			MockJsonPost(c, http.MethodPost, tune)
		})

		When("service returns an error on creation", func() {
			BeforeEach(func() {
				dataService.EXPECT().CreateTune(tune).
					Return(nil, fmt.Errorf("xxx"))
			})

			It("should return a server error", func() {
				Expect(httpRec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("service successfully creates tune", func() {
			BeforeEach(func() {
				dataService.EXPECT().CreateTune(tune).
					Return(&apimodel.Tune{
						ID:    1,
						Title: tune.Title,
					}, nil)
			})

			It("should return ok and the tune", func() {
				Expect(httpRec.Code).To(Equal(http.StatusOK))
				data, err := io.ReadAll(httpRec.Body)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(data)).To(Equal("{\"id\":1,\"title\":\"test title\"}"))
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
