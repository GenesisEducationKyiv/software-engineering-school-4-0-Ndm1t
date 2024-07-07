package services

import (
	"context"
	"informing-service/internal/models"
	"log"
)

const sendEmailEvent = "SendEmail"

type (
	SubscriptionRepositoryInterface interface {
		Create(email string) (*models.Subscription, error)
		ListSubscribed() ([]models.Subscription, error)
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
	subscriptions, err := s.subscriptionRepository.ListSubscribed()
	log.Printf("Subscriptions fetched: %v", subscriptions)
	if err != nil {
		log.Printf("failed to list subscribed emails: %v", err.Error())
		return
	}

	for _, v := range subscriptions {
		err = s.subscriptionProducer.Publish(sendEmailEvent, v, context.Background())
		if err != nil {
			log.Printf("failed to publish SendEmail command: %v", err)
		}
	}

	return
}
