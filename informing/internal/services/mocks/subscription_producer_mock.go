package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"informing-service/internal/models"
)

type MockSubscriptionProducer struct {
	mock.Mock
}

func (m *MockSubscriptionProducer) Publish(eventType string, subscription models.Subscription, ctx context.Context) error {
	args := m.Called(eventType, subscription, ctx)
	return args.Error(0)
}
