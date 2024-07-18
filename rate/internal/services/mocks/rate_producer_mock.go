package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"rate-service/internal/models"
)

type MockRateProducer struct {
	mock.Mock
}

func (m *MockRateProducer) Publish(eventType string, rate models.Rate, ctx context.Context) error {
	args := m.Called(eventType, rate, ctx)
	return args.Error(0)
}
