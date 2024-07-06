package services

import (
	"informing-service/internal/models"
	"log"
)

type (
	IEmailSender interface {
		SendInforming(subscriptions []string, rate float64)
	}

	SubscriptionRepositoryInterface interface {
		Create(email string) (*models.Subscription, error)
		ListSubscribed() ([]models.Subscription, error)
		Update(subscription models.Subscription) (*models.Subscription, error)
	}

	RateRepositoryInterface interface {
		Create(rate models.Rate) (*models.Rate, error)
		GetLatest() (*models.Rate, error)
	}

	InformingServiceInterface interface {
		SendEmails()
	}

	InformingService struct {
		EmailSender            IEmailSender
		subscriptionRepository SubscriptionRepositoryInterface
		rateRepository         RateRepositoryInterface
	}
)

func NewInformingService(subscriptionRepository SubscriptionRepositoryInterface,
	rateRepository RateRepositoryInterface,
	sender IEmailSender) *InformingService {
	return &InformingService{
		EmailSender:            sender,
		rateRepository:         rateRepository,
		subscriptionRepository: subscriptionRepository,
	}
}

func (s *InformingService) SendEmails() {
	rate, err := s.rateRepository.GetLatest()
	log.Printf("Rate fetched: %v", rate)
	if err != nil {
		return
	}

	subscriptions, err := s.subscriptionRepository.ListSubscribed()
	log.Printf("Subscriptions fetched: %v", subscriptions)
	if err != nil {
		log.Printf("failed to list subscribed emails: %v", err.Error())
		return
	}

	var subscribedEmails []string

	for _, v := range subscriptions {
		subscribedEmails = append(subscribedEmails, v.Email)
	}

	s.EmailSender.SendInforming(subscribedEmails, rate.Rate)
	return
}
