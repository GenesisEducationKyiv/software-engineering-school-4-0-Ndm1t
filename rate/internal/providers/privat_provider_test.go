package providers

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"rate-service/internal/app_errors"
	"testing"
)

func TestFetchRateFromPrivat(t *testing.T) {
	t.Run("successful fetch", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			response := [...]GetPrivatResponse{
				{
					Ccy:     "EUR",
					BaseCcy: "UAH",
					Buy:     "43.18000",
					Sell:    "44.18000",
				},
				{
					Ccy:     "USD",
					BaseCcy: "UAH",
					Buy:     "40.30000",
					Sell:    "40.90000",
				},
			}
			_ = json.NewEncoder(w).Encode(response)
		}
		server, teardown := mockServer(handler)
		defer teardown()
		restoreEnv := setEnv("PRIVAT_URL", server.URL)
		defer restoreEnv()

		provider := NewPrivatProvider()
		rate, err := provider.FetchRate()
		assert.NoError(t, err)
		assert.NotNil(t, rate)
		assert.Equal(t, 40.9, *rate)
	})

	t.Run("http request error", func(t *testing.T) {
		restoreEnv := setEnv("PRIVAT_URL", "http://invalid-url")
		defer restoreEnv()

		provider := NewPrivatProvider()
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
		restoreEnv := setEnv("PRIVAT_URL", server.URL)
		defer restoreEnv()

		provider := NewPrivatProvider()
		rate, err := provider.FetchRate()
		assert.Error(t, err)
		assert.Nil(t, rate)
		assert.True(t, errors.Is(err, apperrors.ErrRateFetch))
	})

}
