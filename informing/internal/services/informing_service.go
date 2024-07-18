package services

import (
	"context"
	"informing-service/internal/models"
	"log"
	"sync"
)

const sendEmailEvent = "SendEmail"

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
	}
)

func NewInformingService(subscriptionRepository SubscriptionRepositoryInterface,
	subscriptionProducer SubscriptionProducerInterface) *InformingService {
	return &InformingService{
		subscriptionRepository: subscriptionRepository,
		subscriptionProducer:   subscriptionProducer,
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
		log.Printf("Subscriptions fetched: %v", subscriptions)
		if err != nil {
			log.Printf("failed to list subscribed emails: %v", err.Error())
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
				log.Printf("Spawned goroutine")
				err = s.subscriptionProducer.Publish(sendEmailEvent, subscription, context.Background())
				if err != nil {
					log.Printf("failed to publish SendEmail command: %v", err)
				}
				wg.Done()
			}(subscription)
		}
		wg.Wait()
	}

	return
}
