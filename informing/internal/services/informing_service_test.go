package services

import (
	"context"
	"errors"
	"informing-service/internal/models"
	"informing-service/internal/services/mocks"
	"testing"
)

func TestInformingService_SendEmails(t *testing.T) {
	mockSubscriptionRepo := new(mocks.MockSubscriptionRepository)
	mockSubscriptionProducer := new(mocks.MockSubscriptionProducer)

	service := NewInformingService(mockSubscriptionRepo, mockSubscriptionProducer)

	t.Run("success", func(t *testing.T) {
		mockSubscriptions := []models.Subscription{
			{Email: "test1@example.com"},
			{Email: "test2@example.com"},
		}

		mockSubscriptionRepo.On("ListSubscribed").Return(mockSubscriptions, nil)
		mockSubscriptionProducer.On("Publish", "SendEmail", mockSubscriptions[0], context.Background()).
			Return(nil)
		mockSubscriptionProducer.On("Publish", "SendEmail", mockSubscriptions[1], context.Background()).
			Return(nil)

		service.SendEmails()

		mockSubscriptionRepo.AssertCalled(t, "ListSubscribed")
		mockSubscriptionProducer.AssertCalled(t, "Publish", "SendEmail", mockSubscriptions[0], context.Background())
		mockSubscriptionProducer.AssertCalled(t, "Publish", "SendEmail", mockSubscriptions[1], context.Background())
	})

	t.Run("subscription fetch failure", func(t *testing.T) {

		mockSubscriptionRepo.On("ListSubscribed").Return(nil, errors.New("subscription fetch error"))
		service.SendEmails()

		mockSubscriptionRepo.AssertCalled(t, "ListSubscribed")
		mockSubscriptionProducer.AssertNotCalled(t, "Publish")
	})
}
