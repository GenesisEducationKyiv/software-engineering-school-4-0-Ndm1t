package providers

import (
	"encoding/json"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"rate-service/internal/app_errors"
)

type GetRateResponse struct {
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

type ExchangeAPIProvider struct{}

func NewExchangeAPIProvider() *ExchangeAPIProvider {
	return &ExchangeAPIProvider{}
}

func (p *ExchangeAPIProvider) FetchRate() (*float64, error) {
	res, err := http.Get(viper.GetString("API_URL"))
	if err != nil {
		return nil, apperrors.ErrRateFetch
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, apperrors.ErrRateFetch
	}

	var rateList GetRateResponse
	err = json.Unmarshal(body, &rateList)
	if err != nil {
		return nil, apperrors.ErrRateFetch
	}

	rate := rateList.ConversionRates["UAH"]

	return &rate, nil
}
