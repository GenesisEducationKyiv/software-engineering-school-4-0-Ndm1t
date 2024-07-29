package services

import (
	"context"
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"subscription-service/internal/app_errors"
	"subscription-service/internal/models"
)

const (
	subscriptionCreatedEvent = "SubscriptionCreated"
	subscriptionDeletedEvent = "SubscriptionDeleted"
	publishedSuccessfully    = "published_successfully"
	publishedFail            = "published_fail"
)

type (
	Logger interface {
		Warnf(template string, arguments ...interface{})
		Infof(template string, arguments ...interface{})
		Warn(arguments ...interface{})
		Info(arguments ...interface{})
	}

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
		logger                   Logger
	}
)

func NewSubscriptionService(
	subscriptionRepository ISubscriptionDao,
	subscriptionProducer SubscriptionProducerInterface,
	subscriptionSagaProducer SubscriptionProducerInterface,
	logger Logger) *SubscriptionService {
	return &SubscriptionService{
		SubscriptionDao:          subscriptionRepository,
		SubscriptionProducer:     subscriptionProducer,
		SubscriptionSagaProducer: subscriptionSagaProducer,
		logger:                   logger,
	}
}

func (s *SubscriptionService) Subscribe(email string) (*models.Email, error) {
	subscription, err := s.SubscriptionDao.Find(email)
	if err != nil {
		s.logger.Warnf("couldn't find subscription: %v", err.Error())
		return nil, apperrors.ErrDatabase
	}

	if subscription != nil && subscription.Status == models.Subscribed {
		s.logger.Infof("subscripiton already exists: %v", subscription.Email)
		return nil, apperrors.ErrSubscriptionAlreadyExists
	}

	if subscription.Status == models.Unsubscribed {
		subscription.Status = models.Subscribed
		subscription, err = s.SubscriptionDao.Update(*subscription)
		if err != nil {
			s.logger.Warnf("failed to update subscription status: %v", err.Error())
			return nil, apperrors.ErrDatabase
		}
	}

	if (models.Email{}) == *subscription {
		subscription, err = s.SubscriptionDao.Create(email)
		if err != nil {
			s.logger.Warnf("failed to create subscription: %v", err.Error())
			return nil, apperrors.ErrDatabase
		}
	}

	err = s.SubscriptionSagaProducer.Publish(subscriptionCreatedEvent, *subscription, context.Background())
	if err != nil {
		s.logger.Warnf("failed to publish message: %v", err.Error())
		metrics.GetOrCreateCounter(fmt.Sprintf(`%v{email=%q}`, publishedFail, subscriptionCreatedEvent)).Inc()
		return nil, err
	}
	metrics.GetOrCreateCounter(fmt.Sprintf(`%v{email=%q}`, publishedSuccessfully, subscriptionCreatedEvent)).Inc()

	return subscription, err
}

func (s *SubscriptionService) ListSubscribed() ([]string, error) {
	subscriptions, err := s.SubscriptionDao.ListSubscribed()
	if err != nil {
		s.logger.Warnf("failed to list subscriptions: %v", err.Error())
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
		s.logger.Warnf("couldn't find subscription: %v", err.Error())
		return apperrors.ErrDatabase
	}

	if *subscription == (models.Email{}) {
		s.logger.Info("subscription does not exist")
		return fmt.Errorf("subscription does not exist")
	}

	if *subscription != (models.Email{}) && subscription.Status == models.Unsubscribed {
		s.logger.Infof("subscription already unsubscribed: %v", subscription.Email)
		return apperrors.ErrAlreadyUnsubscribed
	}

	if subscription.Status == models.Subscribed {
		subscription.Status = models.Unsubscribed
		_, err = s.SubscriptionDao.Update(*subscription)
		if err != nil {
			s.logger.Warnf("failed to update subscription: %v", err.Error())
			return apperrors.ErrDatabase
		}
	}

	err = s.SubscriptionProducer.Publish(subscriptionDeletedEvent, *subscription, context.Background())
	if err != nil {
		metrics.GetOrCreateCounter(fmt.Sprintf(`%v{email=%q}`, publishedFail, subscriptionDeletedEvent)).Inc()
		s.logger.Warn("failed to publish DleteSubscription event", subscription)
		return fmt.Errorf("failed to publish DleteSubscription event")
	}
	metrics.GetOrCreateCounter(fmt.Sprintf(`%v{email=%q}`, publishedSuccessfully, subscriptionDeletedEvent)).Inc()

	return nil
}

func (s *SubscriptionService) Delete(email string) error {
	subscription, err := s.SubscriptionDao.Find(email)
	if err != nil {
		s.logger.Warnf("couldn't find subscription: %v", err.Error())
		return err
	}
	err = s.SubscriptionDao.Delete(subscription)
	if err != nil {
		s.logger.Warnf("failed to delete subscription: %v", err.Error())
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateSate(email string, state models.State) error {
	subscription, err := s.SubscriptionDao.Find(email)
	if err != nil {
		s.logger.Warnf("couldn't find subscription: %v", err.Error())
		return err
	}
	subscription.State = state
	_, err = s.SubscriptionDao.Update(*subscription)
	if err != nil {
		s.logger.Warnf("failed to update subscription: %v", err.Error())
		return err
	}
	err = s.SubscriptionProducer.Publish(subscriptionCreatedEvent, *subscription, context.Background())
	if err != nil {
		metrics.GetOrCreateCounter(fmt.Sprintf(`%v{email=%q}`, publishedFail, subscriptionCreatedEvent)).Inc()
		s.logger.Warnf("failed to publish message: %v", err.Error())
		return err
	}
	metrics.GetOrCreateCounter(fmt.Sprintf(`%v{email=%q}`, publishedSuccessfully, subscriptionCreatedEvent)).Inc()
	return nil
}
