package providers

import (
	"encoding/json"
	"gses4_project/internal/apperrors"
	"io"
	"net/http"
	"os"
)

type GetRateResponse struct {
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

type ExchangeAPIProvider struct{}

func NewExchangeAPIProvider() *ExchangeAPIProvider {
	return &ExchangeAPIProvider{}
}

func (p *ExchangeAPIProvider) FetchRate() (*float64, error) {
	res, err := http.Get(os.Getenv("API_URL"))
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
