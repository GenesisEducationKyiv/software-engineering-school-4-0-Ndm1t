package services

type IRateAPIProvider interface {
	FetchRate() (*float64, error)
}

type IRateService interface {
	Get() (*float64, error)
}

type RateService struct {
	APIProvider IRateAPIProvider
}

func NewRateService(provider IRateAPIProvider) *RateService {
	return &RateService{
		APIProvider: provider,
	}
}

func (r *RateService) Get() (*float64, error) {
	rate, err := r.APIProvider.FetchRate()
	if err != nil {
		return nil, err
	}

	return rate, nil
}
