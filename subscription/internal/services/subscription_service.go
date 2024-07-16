package services

import (
	"context"
	"fmt"
	"subscription-service/internal/app_errors"
	"subscription-service/internal/models"
)

const (
	subscriptionCreatedEvent = "SubscriptionCreated"
	subscriptionDeletedEvent = "SubscriptionDeleted"
)

type (
	ISubscriptionDao interface {
		Find(email string) (*models.Email, error)
		Create(email string) (*models.Email, error)
		ListSubscribed() ([]models.Email, error)
		Update(subscription models.Email) (*models.Email, error)
		Delete(subscription *models.Email) error
	}

	SubscriptionProducerInterface interface {
		Publish(eventType string, email models.Email, ctx context.Context) error
	}

	ISubscriptionService interface {
		Subscribe(email string) (*models.Email, error)
		ListSubscribed() ([]string, error)
		Unsubscribe(email string) error
		UpdateSate(email string, state models.State) error
		Delete(email string) error
	}

	SubscriptionService struct {
		SubscriptionDao          ISubscriptionDao
		SubscriptionProducer     SubscriptionProducerInterface
		SubscriptionSagaProducer SubscriptionProducerInterface
	}
)

func NewSubscriptionService(
	subscriptionRepository ISubscriptionDao,
	subscriptionProducer SubscriptionProducerInterface,
	subscriptionSagaProducer SubscriptionProducerInterface) *SubscriptionService {
	return &SubscriptionService{
		SubscriptionDao:          subscriptionRepository,
		SubscriptionProducer:     subscriptionProducer,
		SubscriptionSagaProducer: subscriptionSagaProducer,
	}
}

func (s *SubscriptionService) Subscribe(email string) (*models.Email, error) {
	subscription, err := s.SubscriptionDao.Find(email)
	if err != nil {
		return nil, apperrors.ErrDatabase
	}

	if subscription != nil && subscription.Status == models.Subscribed {
		return nil, apperrors.ErrSubscriptionAlreadyExists
	}

	if subscription.Status == models.Unsubscribed {
		subscription.Status = models.Subscribed
		subscription, err = s.SubscriptionDao.Update(*subscription)
		if err != nil {
			return nil, apperrors.ErrDatabase
		}
	}

	if (models.Email{}) == *subscription {
		subscription, err = s.SubscriptionDao.Create(email)
		if err != nil {
			return nil, apperrors.ErrDatabase
		}
	}

	err = s.SubscriptionSagaProducer.Publish(subscriptionCreatedEvent, *subscription, context.Background())
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

func (s *SubscriptionService) Unsubscribe(email string) error {
	subscription, err := s.SubscriptionDao.Find(email)
	if err != nil {
		return apperrors.ErrDatabase
	}
	if *subscription == (models.Email{}) {
		return fmt.Errorf("subscription does not exist")
	}

	if *subscription != (models.Email{}) && subscription.Status == models.Unsubscribed {
		return apperrors.ErrAlreadyUnsubscribed
	}

	if subscription.Status == models.Subscribed {
		subscription.Status = models.Unsubscribed
		_, err = s.SubscriptionDao.Update(*subscription)
		if err != nil {
			return apperrors.ErrDatabase
		}
	}

	err = s.SubscriptionProducer.Publish(subscriptionDeletedEvent, *subscription, context.Background())
	if err != nil {
		return fmt.Errorf("failed to publish DleteSubscription event")
	}

	return nil
}

func (s *SubscriptionService) Delete(email string) error {
	subscription, err := s.SubscriptionDao.Find(email)
	if err != nil {
		return err
	}
	err = s.SubscriptionDao.Delete(subscription)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateSate(email string, state models.State) error {
	subscription, err := s.SubscriptionDao.Find(email)
	if err != nil {
		return err
	}
	subscription.State = state
	_, err = s.SubscriptionDao.Update(*subscription)
	if err != nil {
		return err
	}
	err = s.SubscriptionProducer.Publish(subscriptionCreatedEvent, *subscription, context.Background())
	if err != nil {
		return err
	}
	return nil
}
