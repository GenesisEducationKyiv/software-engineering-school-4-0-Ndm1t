package providers

import (
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"rate-service/internal/app_errors"
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
	oldValue := viper.GetString(key)
	viper.Set(key, value)
	return func() {
		viper.Set(key, oldValue)
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

		provider := NewExchangeAPIProvider()
		rate, err := provider.FetchRate()
		assert.NoError(t, err)
		assert.NotNil(t, rate)
		assert.Equal(t, 27.5, *rate)
	})

	t.Run("http request error", func(t *testing.T) {
		restoreEnv := setEnv("API_URL", "http://invalid-url")
		defer restoreEnv()

		provider := NewExchangeAPIProvider()
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

		provider := NewExchangeAPIProvider()
		rate, err := provider.FetchRate()
		assert.Error(t, err)
		assert.Nil(t, rate)
		assert.True(t, errors.Is(err, apperrors.ErrRateFetch))
	})

}
