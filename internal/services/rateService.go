package services

import (
	"gses4_project/internal/container"
	"gses4_project/internal/pkg"
)

type IRateAPIProvider interface {
	FetchRate() (*float64, error)
}

type IRateService interface {
	Get() (*float64, error)
}

type RateService struct {
	APIProvider IRateAPIProvider
	container   container.IContainer
}

func NewRateService(container container.IContainer) *RateService {
	return &RateService{
		APIProvider: pkg.NewUSDRateAPIProvider(),
		container:   container,
	}
}

func (r *RateService) Get() (*float64, error) {
	return r.APIProvider.FetchRate()
}
