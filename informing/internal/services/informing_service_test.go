package services

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"informing-service/internal/models"
	"testing"
)

// Mocking the interfaces
type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendInforming(subscriptions []string, rate float64) {
	m.Called(subscriptions, rate)
}

type MockSubscriptionRepository struct {
	mock.Mock
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

type MockRateRepository struct {
	mock.Mock
}

func (m *MockRateRepository) Create(rate models.Rate) (*models.Rate, error) {
	args := m.Called(rate)
	return args.Get(0).(*models.Rate), args.Error(1)
}

func (m *MockRateRepository) GetLatest() (*models.Rate, error) {
	args := m.Called()
	return args.Get(0).(*models.Rate), args.Error(1)
}

func TestInformingService_SendEmails(t *testing.T) {
	mockEmailSender := new(MockEmailSender)
	mockSubscriptionRepo := new(MockSubscriptionRepository)
	mockRateRepo := new(MockRateRepository)

	service := NewInformingService(mockSubscriptionRepo, mockRateRepo, mockEmailSender)

	t.Run("success", func(t *testing.T) {
		mockRate := &models.Rate{Rate: 42.0}
		mockSubscriptions := []models.Subscription{
			{Email: "test1@example.com"},
			{Email: "test2@example.com"},
		}
		subscribedEmails := []string{"test1@example.com", "test2@example.com"}

		mockRateRepo.On("GetLatest").Return(mockRate, nil)
		mockSubscriptionRepo.On("ListSubscribed").Return(mockSubscriptions, nil)
		mockEmailSender.On("SendInforming", subscribedEmails, mockRate.Rate).Return(nil)

		service.SendEmails()

		mockRateRepo.AssertCalled(t, "GetLatest")
		mockSubscriptionRepo.AssertCalled(t, "ListSubscribed")
		mockEmailSender.AssertCalled(t, "SendInforming", subscribedEmails, mockRate.Rate)
	})

	t.Run("rate fetch failure", func(t *testing.T) {
		mockRateRepo.On("GetLatest").Return(nil, errors.New("rate fetch error"))

		service.SendEmails()

		mockRateRepo.AssertCalled(t, "GetLatest")
		mockSubscriptionRepo.AssertNotCalled(t, "ListSubscribed")
		mockEmailSender.AssertNotCalled(t, "SendInforming", mock.Anything, mock.Anything)
	})

	t.Run("subscription fetch failure", func(t *testing.T) {
		mockRate := &models.Rate{Rate: 42.0}

		mockRateRepo.On("GetLatest").Return(mockRate, nil)
		mockSubscriptionRepo.On("ListSubscribed").Return(nil, errors.New("subscription fetch error"))

		service.SendEmails()

		mockRateRepo.AssertCalled(t, "GetLatest")
		mockSubscriptionRepo.AssertCalled(t, "ListSubscribed")
		mockEmailSender.AssertNotCalled(t, "SendInforming", mock.Anything, mockRate.Rate)
	})
}
