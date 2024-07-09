package services

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"subscription-service/internal/app_errors"
	"subscription-service/internal/models"
	"testing"
)

type MockSubscriptionDao struct {
	mock.Mock
}

func (m *MockSubscriptionDao) Find(email string) (*models.Email, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Email), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSubscriptionDao) Create(email string) (*models.Email, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Email), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSubscriptionDao) ListSubscribed() ([]models.Email, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]models.Email), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSubscriptionDao) Update(subscription models.Email) (*models.Email, error) {
	args := m.Called(subscription)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Email), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestSubscriptionService_Subscribe_AlreadySubscribed(t *testing.T) {
	mockDao := new(MockSubscriptionDao)
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
	mockDao := new(MockSubscriptionDao)
	email := "test@example.com"
	subscription := &models.Email{Email: email, Status: models.Unsubscribed}
	updatedSubscription := &models.Email{Email: email, Status: models.Subscribed}
	mockDao.On("Find", email).Return(subscription, nil)
	mockDao.On("Update", *updatedSubscription).Return(updatedSubscription, nil)

	service := &SubscriptionService{
		SubscriptionDao: mockDao,
	}

	result, err := service.Subscribe(email)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedSubscription, result)

	mockDao.AssertExpectations(t)
}

func TestSubscriptionService_Subscribe_NotFound(t *testing.T) {
	mockDao := new(MockSubscriptionDao)
	email := "test@example.com"
	newSubscription := &models.Email{Email: email, Status: models.Subscribed}
	mockDao.On("Find", email).Return(&models.Email{}, nil)
	mockDao.On("Create", email).Return(newSubscription, nil)

	service := &SubscriptionService{
		SubscriptionDao: mockDao,
	}

	result, err := service.Subscribe(email)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newSubscription, result)

	mockDao.AssertExpectations(t)
}
