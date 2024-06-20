package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gses4_project/internal/apperrors"
	"gses4_project/internal/container"
	"gses4_project/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockRateService struct {
	mock.Mock
}

func (m *MockRateService) Get() (*float64, error) {
	args := m.Called()
	if args.Get(0) != nil {
		rate := args.Get(0).(float64)
		return &rate, args.Error(1)
	}
	return nil, args.Error(1)
}

func setupTestServer(rateService services.IRateService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	container := &container.Container{}
	controller := &RateController{
		RateService: rateService,
		container:   container,
	}
	router.GET("/rate", controller.Get)
	return router
}

func TestRateController_Get_Success(t *testing.T) {
	mockRateService := new(MockRateService)
	rate := 27.5
	mockRateService.On("Get").Return(rate, nil)

	router := setupTestServer(mockRateService)

	req, _ := http.NewRequest(http.MethodGet, "/rate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `27.5`, w.Body.String())
}

func TestRateController_Get_ServiceError(t *testing.T) {
	mockRateService := new(MockRateService)
	mockRateService.On("Get").Return(nil, apperrors.ErrRateFetch)

	router := setupTestServer(mockRateService)

	req, _ := http.NewRequest(http.MethodGet, "/rate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, apperrors.ErrRateFetch.StatusCode, w.Code)
	assert.JSONEq(t, `{"error":"failed to fetch rate"}`, w.Body.String())
}

func TestRateController_Get_InternalServerError(t *testing.T) {
	mockRateService := new(MockRateService)
	mockRateService.On("Get").Return(nil, errors.New("unexpected error"))

	router := setupTestServer(mockRateService)

	req, _ := http.NewRequest(http.MethodGet, "/rate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, apperrors.ErrInternalServer.StatusCode, w.Code)
	assert.JSONEq(t, `{"error":"internal server error"}`, w.Body.String())
}
