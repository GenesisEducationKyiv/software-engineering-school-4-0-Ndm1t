package services

import (
	"gses4_project/internal/models"
	"gses4_project/internal/pkg"
	"log"
)

type IEmailSender interface {
	SendInforming(subscriptions []models.Email, rate float64)
}

type InformingService struct {
	EmailSender         IEmailSender
	SubscriptionService *SubscriptionService
	RateService         *RateService
}

func NewInformingService() *InformingService {
	return &InformingService{
		EmailSender:         pkg.NewSmtpEmailSender(),
		SubscriptionService: NewSubscriptionService(),
		RateService:         NewRateService(),
	}
}

func (s *InformingService) SendEmails() {
	rate, err := s.RateService.Get()
	if err != nil {
		log.Printf("Failed to fetch rate: %v", err.Error())
		return
	}

	subscriptions, err := s.SubscriptionService.SubscriptionDao.ListSubscribed()
	if err != nil {
		log.Printf("Failed to fetch rate: %v", err.Error())
		return
	}

	s.EmailSender.SendInforming(subscriptions, *rate)
}
