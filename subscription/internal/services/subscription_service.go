package services

import (
	"context"
	"subscription-service/internal/app_errors"
	"subscription-service/internal/models"
)

const subscriptionCreatedEvent = "SubscriptionCreated"

type (
	ISubscriptionDao interface {
		Find(email string) (*models.Email, error)
		Create(email string) (*models.Email, error)
		ListSubscribed() ([]models.Email, error)
		Update(subscription models.Email) (*models.Email, error)
	}

	SubscriptionProducerInterface interface {
		Publish(eventType string, email models.Email, ctx context.Context) error
	}

	ISubscriptionService interface {
		Subscribe(email string) (*models.Email, error)
		ListSubscribed() ([]string, error)
	}

	SubscriptionService struct {
		SubscriptionDao      ISubscriptionDao
		SubscriptionProducer SubscriptionProducerInterface
	}
)

func NewSubscriptionService(
	subscriptionRepository ISubscriptionDao,
	subscriptionProducer SubscriptionProducerInterface) *SubscriptionService {
	return &SubscriptionService{
		SubscriptionDao:      subscriptionRepository,
		SubscriptionProducer: subscriptionProducer,
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

	err = s.SubscriptionProducer.Publish(subscriptionCreatedEvent, *subscription, context.Background())
	if err != nil {
		return nil, err
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
