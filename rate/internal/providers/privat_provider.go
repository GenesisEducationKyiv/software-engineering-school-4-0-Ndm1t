package providers

import (
	"encoding/json"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"rate-service/internal/app_errors"
	"strconv"
)

const (
	USD = "USD"
)

type (
	GetPrivatResponse struct {
		Ccy     string `json:"ccy"`
		BaseCcy string `json:"base_ccy"`
		Buy     string `json:"buy"`
		Sell    string `json:"sale"`
	}

	PrivatProvider struct{}
)

func NewPrivatProvider() *PrivatProvider {
	return &PrivatProvider{}
}

func (p *PrivatProvider) FetchRate() (*float64, error) {
	res, err := http.Get(viper.GetString("PRIVAT_URL"))
	if err != nil {
		return nil, apperrors.ErrRateFetch
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, apperrors.ErrRateFetch
	}

	var rateList []GetPrivatResponse
	err = json.Unmarshal(body, &rateList)
	if err != nil {
		return nil, apperrors.ErrRateFetch
	}

	var rate float64
	for _, v := range rateList {
		if v.Ccy == USD {
			rate, err = strconv.ParseFloat(v.Sell, 64)
			if err != nil {
				return nil, err
			}
		}
	}

	return &rate, nil
}
