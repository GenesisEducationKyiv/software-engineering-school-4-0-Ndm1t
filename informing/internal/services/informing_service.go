package services

import (
	"log"
)

type (
	IEmailSender interface {
		SendInforming(subscriptions []string, rate float64)
	}
	SubscriptionClientInterface interface {
		ListSubscribed() ([]string, error)
	}
	RateClientInterface interface {
		FetchRate() (*float64, error)
	}
	InformingServiceInterface interface {
		SendEmails()
	}
	InformingService struct {
		EmailSender        IEmailSender
		subscriptionClient SubscriptionClientInterface
		rateClient         RateClientInterface
	}
)

func NewInformingService(subscriptionClient SubscriptionClientInterface, rateClient RateClientInterface, sender IEmailSender) *InformingService {
	return &InformingService{
		EmailSender:        sender,
		subscriptionClient: subscriptionClient,
		rateClient:         rateClient,
	}
}

func (s *InformingService) SendEmails() {
	rate, err := s.rateClient.FetchRate()
	log.Printf("Rate fetched: %v", rate)
	if err != nil {
		return
	}

	subscriptions, err := s.subscriptionClient.ListSubscribed()
	log.Printf("Subscriptions fetched: %v", subscriptions)
	if err != nil {
		log.Printf("failed to list subscribed emails: %v", err.Error())
		return
	}
	s.EmailSender.SendInforming(subscriptions, *rate)
	return
}
