package pkg

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

type USDRateAPIProvider struct{}

func NewUSDRateAPIProvider() *USDRateAPIProvider {
	return &USDRateAPIProvider{}
}

func (p *USDRateAPIProvider) FetchRate() (*float64, error) {
	res, err := http.Get(os.Getenv("API_URL"))
	if err != nil {
		return nil, apperrors.NewHttpError("Failed to fetch rate from API",
			http.StatusInternalServerError)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, apperrors.NewHttpError("Failed to read response from rate API",
			http.StatusInternalServerError)
	}

	var rateList GetRateResponse
	err = json.Unmarshal(body, &rateList)
	if err != nil {
		return nil, apperrors.NewHttpError("Failed to read response from rate API",
			http.StatusInternalServerError)
	}

	rate := rateList.ConversionRates["UAH"]

	return &rate, nil
}
