package services

import (
	"github.com/stretchr/testify/mock"
	"gses4_project/internal/apperrors"
	"gses4_project/internal/models"
	"testing"
)

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendInforming(subscriptions []models.Email, rate float64) {
	m.Called(subscriptions, rate)
}

type MockRateService struct {
	mock.Mock
}

func (m *MockRateService) Get() (*float64, error) {
	args := m.Called()
	if args.Get(0) != nil {
		rate := args.Get(0).(float64)
		return &rate, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestInformingService_SendEmails_Success(t *testing.T) {
	mockEmailSender := new(MockEmailSender)
	mockSubscriptionDao := new(MockSubscriptionDao)
	mockRateService := new(MockRateService)

	rate := 27.5
	subscriptions := []models.Email{
		{Email: "test1@example.com", Status: models.Subscribed},
		{Email: "test2@example.com", Status: models.Subscribed},
	}

	mockRateService.On("Get").Return(rate, nil)
	mockSubscriptionDao.On("ListSubscribed").Return(subscriptions, nil)
	mockEmailSender.On("SendInforming", subscriptions, rate).Return()

	service := &InformingService{
		EmailSender:     mockEmailSender,
		SubscriptionDao: mockSubscriptionDao,
		RateService:     mockRateService,
	}

	service.SendEmails()

	mockRateService.AssertExpectations(t)
	mockSubscriptionDao.AssertExpectations(t)
	mockEmailSender.AssertExpectations(t)
}

func TestInformingService_SendEmails_FailRateFetch(t *testing.T) {
	mockEmailSender := new(MockEmailSender)
	mockSubscriptionDao := new(MockSubscriptionDao)
	mockRateService := new(MockRateService)

	mockRateService.On("Get").Return(nil, apperrors.ErrRateFetch)

	service := &InformingService{
		EmailSender:     mockEmailSender,
		SubscriptionDao: mockSubscriptionDao,
		RateService:     mockRateService,
	}

	service.SendEmails()

	mockRateService.AssertExpectations(t)
	mockSubscriptionDao.AssertNotCalled(t, "ListSubscribed")
	mockEmailSender.AssertNotCalled(t, "SendInforming")
}

func TestInformingService_SendEmails_FailSubscriptionFetch(t *testing.T) {
	mockEmailSender := new(MockEmailSender)
	mockSubscriptionDao := new(MockSubscriptionDao)
	mockRateService := new(MockRateService)

	rate := 27.5

	mockRateService.On("Get").Return(rate, nil)
	mockSubscriptionDao.On("ListSubscribed").Return(nil, apperrors.ErrDatabase)

	service := &InformingService{
		EmailSender:     mockEmailSender,
		SubscriptionDao: mockSubscriptionDao,
		RateService:     mockRateService,
	}

	service.SendEmails()

	mockRateService.AssertExpectations(t)
	mockSubscriptionDao.AssertExpectations(t)
	mockEmailSender.AssertNotCalled(t, "SendInforming")
}
