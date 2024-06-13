package services

import "gses4_project/internal/pkg"

type IRateAPIProvider interface {
	FetchRate() (*float64, error)
}

type RateService struct {
	APIProvider IRateAPIProvider
}

func NewRateService() *RateService {
	return &RateService{
		APIProvider: pkg.NewUSDRateAPIProvider(),
	}
}

func (r *RateService) Get() (*float64, error) {
	return r.APIProvider.FetchRate()
}
