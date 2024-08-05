package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"subscription-service/internal/app_errors"
	"subscription-service/internal/models"
	"subscription-service/internal/services/mocks"
	"testing"
)

func TestSubscriptionService_Subscribe_AlreadySubscribed(t *testing.T) {
	mockDao := new(mocks.MockSubscriptionDao)
	email := "test@example.com"
	subscription := &models.Email{Email: email, Status: models.Subscribed}
	mockDao.On("Find", email).Return(subscription, nil)

	service := &SubscriptionService{
		SubscriptionDao: mockDao,
	}

	result, err := service.Subscribe(email)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, apperrors.ErrSubscriptionAlreadyExists, err)

	mockDao.AssertExpectations(t)
}

func TestSubscriptionService_Subscribe_Unsubscribed(t *testing.T) {
	mockDao := new(mocks.MockSubscriptionDao)
	mockProducer := new(mocks.MockSubscriptionProducer)
	email := "test@example.com"
	subscription := &models.Email{Email: email, Status: models.Unsubscribed}
	updatedSubscription := &models.Email{Email: email, Status: models.Subscribed}
	mockDao.On("Find", email).Return(subscription, nil)
	mockDao.On("Update", *updatedSubscription).Return(updatedSubscription, nil)
	mockProducer.On("Publish", "SubscriptionCreated", *updatedSubscription, context.Background()).
		Return(nil)

	service := &SubscriptionService{
		SubscriptionDao:      mockDao,
		SubscriptionProducer: mockProducer,
	}

	result, err := service.Subscribe(email)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedSubscription, result)

	mockDao.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestSubscriptionService_Subscribe_Success(t *testing.T) {
	mockDao := new(mocks.MockSubscriptionDao)
	mockProducer := new(mocks.MockSubscriptionProducer)
	email := "test@example.com"
	newSubscription := &models.Email{Email: email, Status: models.Subscribed}
	mockDao.On("Find", email).Return(&models.Email{}, nil)
	mockDao.On("Create", email).Return(newSubscription, nil)
	mockProducer.On("Publish", "SubscriptionCreated", *newSubscription, context.Background()).
		Return(nil)

	service := &SubscriptionService{
		SubscriptionDao:      mockDao,
		SubscriptionProducer: mockProducer,
	}

	result, err := service.Subscribe(email)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newSubscription, result)

	mockDao.AssertExpectations(t)
}
