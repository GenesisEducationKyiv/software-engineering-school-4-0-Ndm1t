package services

import (
	"context"
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"go.uber.org/zap"
	"informing-service/internal/models"
	"sync"
)

const (
	sendEmailEvent        = "SendEmail"
	messagePublishSuccess = "message_publish_success"
	messagePublishFail    = "message_publish_fail"
)

type (
	SubscriptionRepositoryInterface interface {
		Create(email string) (*models.Subscription, error)
		ListSubscribed(limit int, email string) ([]models.Subscription, error)
		Update(subscription models.Subscription) (*models.Subscription, error)
		Delete(subscription models.Subscription) error
	}

	SubscriptionProducerInterface interface {
		Publish(eventType string, subscription models.Subscription, ctx context.Context) error
	}

	InformingServiceInterface interface {
		SendEmails()
	}

	InformingService struct {
		subscriptionRepository SubscriptionRepositoryInterface
		subscriptionProducer   SubscriptionProducerInterface
		logger                 *zap.SugaredLogger
	}
)

func NewInformingService(subscriptionRepository SubscriptionRepositoryInterface,
	subscriptionProducer SubscriptionProducerInterface,
	logger *zap.SugaredLogger) *InformingService {
	return &InformingService{
		subscriptionRepository: subscriptionRepository,
		subscriptionProducer:   subscriptionProducer,
		logger:                 logger,
	}
}

func (s *InformingService) SendEmails() {
	var (
		left       = true
		startEmail = ""
		limit      = 5
	)

	for left {
		subscriptions, err := s.subscriptionRepository.ListSubscribed(limit, startEmail)
		s.logger.Infof("Subscriptions fetched: %v", subscriptions)
		if err != nil {
			s.logger.Warnf("failed to list subscribed emails: %v", err.Error())
			return
		}

		if len(subscriptions) == 0 {
			left = false
		}

		if len(subscriptions) > 0 {
			startEmail = subscriptions[len(subscriptions)-1].Email
		}

		wg := sync.WaitGroup{}
		for _, v := range subscriptions {
			wg.Add(1)
			subscription := v
			go func(subscription models.Subscription) {
				s.logger.Infof("Spawned goroutine")
				err = s.subscriptionProducer.Publish(sendEmailEvent, subscription, context.Background())
				if err != nil {
					metrics.GetOrCreateCounter(fmt.Sprintf(`%v{event=%q}`, messagePublishFail, sendEmailEvent)).Inc()
					s.logger.Warnf("failed to publish SendEmail command: %v", err)
				}
				metrics.GetOrCreateCounter(fmt.Sprintf(`%v{event=%q}`, messagePublishSuccess, sendEmailEvent)).Inc()
				wg.Done()
			}(subscription)
		}
		wg.Wait()
	}

	return
}
