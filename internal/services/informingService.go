package services

import (
	"gses4_project/internal/container"
	"gses4_project/internal/database"
	"gses4_project/internal/models"
	"gses4_project/internal/pkg"
	"log"
)

type IEmailSender interface {
	SendInforming(subscriptions []models.Email, rate float64)
}

type InformingService struct {
	EmailSender     IEmailSender
	SubscriptionDao ISubscriptionDao
	RateService     IRateService
	container       container.IContainer
}

func NewInformingService(container container.IContainer) *InformingService {
	return &InformingService{
		EmailSender:     pkg.NewSmtpEmailSender(),
		SubscriptionDao: database.NewSubcscriptionDao(container.GetDatabase()),
		RateService:     NewRateService(container),
		container:       container,
	}
}

func (s *InformingService) SendEmails() {
	rate, err := s.RateService.Get()
	if err != nil {
		log.Printf("Failed to fetch rate: %v", err.Error())
		return
	}

	subscriptions, err := s.SubscriptionDao.ListSubscribed()
	if err != nil {
		log.Printf("Failed to fetch rate: %v", err.Error())
		return
	}

	s.EmailSender.SendInforming(subscriptions, *rate)
}
