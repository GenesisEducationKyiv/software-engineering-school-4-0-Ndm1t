package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"subscription-service/internal/models"
)

type (
	MockSubscriptionDao struct {
		mock.Mock
	}

	MockSubscriptionProducer struct {
		mock.Mock
	}
)

func (m *MockSubscriptionProducer) Publish(eventType string, email models.Email, ctx context.Context) error {
	args := m.Called(eventType, email, ctx)
	return args.Error(0)
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

func (m *MockSubscriptionDao) Delete(subscription *models.Email) error {
	args := m.Called(subscription)
	if args.Get(0) != nil {
		return args.Error(1)
	}
	return args.Error(1)
}
