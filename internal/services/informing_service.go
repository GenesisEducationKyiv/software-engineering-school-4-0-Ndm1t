package services

import (
	"gses4_project/internal/container"
	"gses4_project/internal/models"
	"log"
)

type IEmailSender interface {
	SendInforming(subscriptions []models.Email, rate float64)
}

type IInformingService interface {
	SendEmails()
}

type InformingService struct {
	EmailSender     IEmailSender
	SubscriptionDao ISubscriptionDao
	RateService     IRateService
	container       container.IContainer
}

func NewInformingService(container container.IContainer, sender IEmailSender,
	subscriptionRepository ISubscriptionDao, rateService IRateService) *InformingService {
	return &InformingService{
		EmailSender:     sender,
		SubscriptionDao: subscriptionRepository,
		RateService:     rateService,
		container:       container,
	}
}

func (s *InformingService) SendEmails() {
	rate, err := s.RateService.Get()
	if err != nil {
		return
	}

	subscriptions, err := s.SubscriptionDao.ListSubscribed()
	if err != nil {
		log.Printf("failed to list subscribed emails: %v", err.Error())
		return
	}

	s.EmailSender.SendInforming(subscriptions, *rate)
}
