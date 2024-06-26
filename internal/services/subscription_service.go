package services

import (
	"gses4_project/internal/apperrors"
	"gses4_project/internal/container"
	"gses4_project/internal/models"
)

type ISubscriptionDao interface {
	Find(email string) (*models.Email, error)
	Create(email string) (*models.Email, error)
	ListSubscribed() ([]models.Email, error)
	Update(subscription models.Email) (*models.Email, error)
}

type ISubscriptionService interface {
	Subscribe(email string) (*models.Email, error)
}

type SubscriptionService struct {
	SubscriptionDao ISubscriptionDao
	container       container.IContainer
}

func NewSubscriptionService(container container.IContainer,
	subscriptionRepository ISubscriptionDao) *SubscriptionService {
	return &SubscriptionService{
		SubscriptionDao: subscriptionRepository,
		container:       container,
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
