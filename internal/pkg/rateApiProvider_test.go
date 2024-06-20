package pkg

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"gses4_project/internal/apperrors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func mockServer(handler http.HandlerFunc) (*httptest.Server, func()) {
	server := httptest.NewServer(handler)
	teardown := func() {
		server.Close()
	}
	return server, teardown
}

func setEnv(key, value string) func() {
	oldValue := os.Getenv(key)
	_ = os.Setenv(key, value)
	return func() {
		_ = os.Setenv(key, oldValue)
	}
}

func TestFetchRate(t *testing.T) {
	t.Run("successful fetch", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			response := GetRateResponse{
				ConversionRates: map[string]float64{
					"UAH": 27.5,
				},
			}
			_ = json.NewEncoder(w).Encode(response)
		}
		server, teardown := mockServer(handler)
		defer teardown()
		restoreEnv := setEnv("API_URL", server.URL)
		defer restoreEnv()

		provider := NewUSDRateAPIProvider()
		rate, err := provider.FetchRate()
		assert.NoError(t, err)
		assert.NotNil(t, rate)
		assert.Equal(t, 27.5, *rate)
	})

	t.Run("http request error", func(t *testing.T) {
		restoreEnv := setEnv("API_URL", "http://invalid-url")
		defer restoreEnv()

		provider := NewUSDRateAPIProvider()
		rate, err := provider.FetchRate()
		assert.Error(t, err)
		assert.Nil(t, rate)
		assert.True(t, errors.Is(err, apperrors.ErrRateFetch))
	})

	t.Run("invalid json response", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`invalid json`))
		}
		server, teardown := mockServer(handler)
		defer teardown()
		restoreEnv := setEnv("API_URL", server.URL)
		defer restoreEnv()

		provider := NewUSDRateAPIProvider()
		rate, err := provider.FetchRate()
		assert.Error(t, err)
		assert.Nil(t, rate)
		assert.True(t, errors.Is(err, apperrors.ErrRateFetch))
	})

}
