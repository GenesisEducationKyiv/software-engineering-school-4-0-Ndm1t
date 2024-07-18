package mocks

import (
	"github.com/stretchr/testify/mock"
	"subscription-service/internal/models"
)

type MockSubscriptionService struct {
	mock.Mock
}

func (m *MockSubscriptionService) Unsubscribe(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockSubscriptionService) Subscribe(email string) (*models.Email, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		email := args.Get(0).(models.Email)
		return &email, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSubscriptionService) ListSubscribed() ([]string, error) {
	args := m.Called()
	if args.Get(0) != nil {
		emails := args.Get(0).([]string)
		return emails, args.Error(1)
	}
	return nil, args.Error(1)
}
