package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"subscription-service/internal/app_errors"
	"subscription-service/internal/models"
	"subscription-service/internal/services"
	"testing"
)

type MockSubscriptionService struct {
	mock.Mock
}

func (m *MockSubscriptionService) Subscribe(email string) (*models.Email, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		email := args.Get(0).(models.Email)
		return &email, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSubscriptionService) ListSubscribed() ([]string, error) {
	args := m.Called()
	if args.Get(0) != nil {
		emails := args.Get(0).([]string)
		return emails, args.Error(1)
	}
	return nil, args.Error(1)
}

func setupSubscriptionTestServer(subscriptionService services.ISubscriptionService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	controller := &SubscriptionController{
		SubscriptionService: subscriptionService,
	}
	router.POST("/subscribe", controller.Subscribe)
	return router
}

func TestSubscriptionController_Subscribe_Success(t *testing.T) {
	mockService := new(MockSubscriptionService)
	email := "test@example.com"
	subscription := models.Email{Email: email, Status: models.Subscribed}
	mockService.On("Subscribe", email).Return(subscription, nil)

	router := setupSubscriptionTestServer(mockService)

	body, _ := json.Marshal(models.Email{Email: email})
	req, _ := http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody, _ := json.Marshal(subscription)
	assert.JSONEq(t, string(expectedBody), w.Body.String())
}

func TestSubscriptionController_Subscribe_BadRequest(t *testing.T) {
	mockService := new(MockSubscriptionService)

	router := setupSubscriptionTestServer(mockService)

	req, _ := http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"Couldn't bind request"}`, w.Body.String())
}

func TestSubscriptionController_Subscribe_ServiceError(t *testing.T) {
	mockService := new(MockSubscriptionService)
	email := "test@example.com"
	mockService.On("Subscribe", email).Return(nil, apperrors.ErrDatabase)

	router := setupSubscriptionTestServer(mockService)

	body, _ := json.Marshal(models.Email{Email: email})
	req, _ := http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, apperrors.ErrDatabase.StatusCode, w.Code)
	expectedBody, _ := json.Marshal(apperrors.ErrDatabase.JSONResponse)
	assert.JSONEq(t, string(expectedBody), w.Body.String())
}

func TestSubscriptionController_Subscribe_InternalServerError(t *testing.T) {
	mockService := new(MockSubscriptionService)
	email := "test@example.com"
	mockService.On("Subscribe", email).Return(nil, errors.New("unexpected error"))

	router := setupSubscriptionTestServer(mockService)

	body, _ := json.Marshal(models.Email{Email: email})
	req, _ := http.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, apperrors.ErrInternalServer.StatusCode, w.Code)
	expectedBody, _ := json.Marshal(apperrors.ErrInternalServer.JSONResponse)
	assert.JSONEq(t, string(expectedBody), w.Body.String())
}
