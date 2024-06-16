package server

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gses4_project/internal/container"
	"gses4_project/internal/models"
	"gses4_project/internal/pkg"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupSQLiteDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	_ = db.AutoMigrate(&models.Email{})

	return db
}

func setupTestServer() *Server {

	db := setupSQLiteDB()

	pkg.LoadConfig("../../.env")

	cont := container.NewContainer(db)

	testServer := NewServer(cont)
	return testServer
}

func TestGetRate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testServer := setupTestServer()

	// Create a request to send to the endpoint
	req, _ := http.NewRequest("GET", "/api/rate", nil)
	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()
	// Perform the request
	testServer.router.ServeHTTP(w, req)

	// Check if the status code is what you expect
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body
	var response float64
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, response)
}

func TestSubscribe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testServer := setupTestServer()

	t.Run("create new subscription success", func(t *testing.T) {
		payload := map[string]string{"email": "test@example.com"}
		jsonPayload, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/api/subscribe", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		testServer.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "email")
	})
	t.Run("create new subscription error", func(t *testing.T) {
		payload := map[string]string{"email": "test@example.com"}
		jsonPayload, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/api/subscribe", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		testServer.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
	})

}
