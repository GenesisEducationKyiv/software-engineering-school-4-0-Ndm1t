package services

import (
	"subscription-service/internal/apperrors"
	"subscription-service/internal/models"
)

type ISubscriptionDao interface {
	Find(email string) (*models.Email, error)
	Create(email string) (*models.Email, error)
	ListSubscribed() ([]models.Email, error)
	Update(subscription models.Email) (*models.Email, error)
}

type ISubscriptionService interface {
	Subscribe(email string) (*models.Email, error)
	ListSubscribed() ([]string, error)
}

type SubscriptionService struct {
	SubscriptionDao ISubscriptionDao
}

func NewSubscriptionService(
	subscriptionRepository ISubscriptionDao) *SubscriptionService {
	return &SubscriptionService{
		SubscriptionDao: subscriptionRepository,
	}
}

func (s *SubscriptionService) Subscribe(email string) (*models.Email, error) {
	subscription, err := s.SubscriptionDao.Find(email)

	if subscription != nil && subscription.Status == models.Subscribed {
		return nil, apperrors.ErrSubscriptionAlreadyExists
	}

	if subscription.Status == models.Unsubscribed {
		subscription.Status = models.Subscribed
		subscription, err = s.SubscriptionDao.Update(*subscription)
	}

	if (models.Email{}) == *subscription {
		subscription, err = s.SubscriptionDao.Create(email)
	}

	return subscription, err
}

func (s *SubscriptionService) ListSubscribed() ([]string, error) {
	subscriptions, err := s.SubscriptionDao.ListSubscribed()
	if err != nil {
		return nil, apperrors.ErrDatabase
	}

	subscribedEmails := make([]string, 0)

	for _, v := range subscriptions {
		subscribedEmails = append(subscribedEmails, v.Email)
	}

	return subscribedEmails, nil

}
