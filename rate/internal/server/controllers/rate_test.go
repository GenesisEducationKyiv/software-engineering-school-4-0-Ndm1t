package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	apperrors "rate-service/internal/app_errors"
	"rate-service/internal/server/controllers/mocks"
	"rate-service/internal/services"
	"testing"
)

func setupRateTestServer(rateService services.IRateService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	controller := &RateController{
		RateService: rateService,
	}
	router.GET("/rate", controller.Get)
	return router
}

func TestRateController_Get_Success(t *testing.T) {
	mockService := new(mocks.MockRateService)
	rate := 27.5
	mockService.On("Get").Return(rate, nil)

	router := setupRateTestServer(mockService)

	req, _ := http.NewRequest(http.MethodGet, "/rate", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, fmt.Sprintf("%v", rate), w.Body.String())
}

func TestRateController_Get_Failed(t *testing.T) {
	mockService := new(mocks.MockRateService)

	mockService.On("Get").Return(nil, apperrors.ErrRateFetch)

	router := setupRateTestServer(mockService)

	req, _ := http.NewRequest(http.MethodGet, "/rate", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expectedBody, err := json.Marshal(apperrors.ErrRateFetch.JSONResponse)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, string(expectedBody), w.Body.String())
}
