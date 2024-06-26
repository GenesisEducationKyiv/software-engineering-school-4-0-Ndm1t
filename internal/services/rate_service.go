package services

import (
	"gses4_project/internal/container"
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

func NewRateService(provider IRateAPIProvider, container container.IContainer) *RateService {
	return &RateService{
		APIProvider: provider,
		container:   container,
	}
}

func (r *RateService) Get() (*float64, error) {
	rate, err := r.APIProvider.FetchRate()
	if err != nil {
		return nil, err
	}

	return rate, nil
}
