package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"rate-service/internal/app_errors"
	"rate-service/internal/services/mocks"
	"testing"
)

func TestRateService_Get_Success(t *testing.T) {
	mockAPIProvider := new(mocks.MockRateAPIProvider)
	mockRateProducer := new(mocks.MockRateProducer)

	// Mock the FetchRate method
	rate := 27.5

	mockAPIProvider.On("FetchRate").Return(rate, nil)
	mockRateProducer.On("Publish", "RateFetched", mock.Anything, context.Background()).
		Return(nil)

	// Inject the mock provider into the service
	service := &RateService{
		APIProvider:  mockAPIProvider,
		rateProducer: mockRateProducer,
	}

	// Test the Get method
	result, err := service.Get()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, rate, *result)

	// Assert that the mock was called
	mockAPIProvider.AssertExpectations(t)
	mockRateProducer.AssertExpectations(t)
}

func TestRateService_Get_FetchRateError(t *testing.T) {
	mockAPIProvider := new(mocks.MockRateAPIProvider)

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
