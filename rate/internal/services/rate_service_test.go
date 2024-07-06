package services

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"rate-service/internal/app_errors"
	"testing"
)

type MockRateAPIProvider struct {
	mock.Mock
}

func (m *MockRateAPIProvider) FetchRate() (*float64, error) {
	args := m.Called()
	if args.Get(0) != nil {
		rate := args.Get(0).(float64)
		return &rate, args.Error(1)
	}
	return nil, args.Error(1)
}
func TestRateService_Get(t *testing.T) {
	mockAPIProvider := new(MockRateAPIProvider)

	// Mock the FetchRate method
	rate := 27.5
	mockAPIProvider.On("FetchRate").Return(rate, nil)

	// Inject the mock provider into the service
	service := &RateService{
		APIProvider: mockAPIProvider,
	}

	// Test the Get method
	result, err := service.Get()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, rate, *result)

	// Assert that the mock was called
	mockAPIProvider.AssertExpectations(t)
}

func TestRateService_Get_FetchRateError(t *testing.T) {
	mockAPIProvider := new(MockRateAPIProvider)

	// Mock the FetchRate method to return an error
	mockAPIProvider.On("FetchRate").Return(nil, apperrors.ErrRateFetch)

	// Inject the mock provider into the service
	service := &RateService{
		APIProvider: mockAPIProvider,
	}

	// Test the Get method
	result, err := service.Get()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, apperrors.ErrRateFetch, err)

	// Assert that the mock was called
	mockAPIProvider.AssertExpectations(t)
}
