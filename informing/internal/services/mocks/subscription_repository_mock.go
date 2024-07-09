package mocks

import (
	"github.com/stretchr/testify/mock"
	"informing-service/internal/models"
)

type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) Delete(subscription models.Subscription) error {
	args := m.Called(subscription)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) Create(email string) (*models.Subscription, error) {
	args := m.Called(email)
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) ListSubscribed() ([]models.Subscription, error) {
	args := m.Called()
	return args.Get(0).([]models.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) Update(subscription models.Subscription) (*models.Subscription, error) {
	args := m.Called(subscription)
	return args.Get(0).(*models.Subscription), args.Error(1)
}
