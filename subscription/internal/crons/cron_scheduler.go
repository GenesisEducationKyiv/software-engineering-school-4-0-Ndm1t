package crons

import (
	"context"
	"github.com/robfig/cron/v3"
	"log"
	"subscription-service/internal/kafka/producers"
	"subscription-service/internal/services"
)

const (
	EveryDayAt9AM = "* * * * *"
)

type CronScheduler struct {
	Cron                *cron.Cron
	subscriptionService services.ISubscriptionService
	emailProducer       *producers.EmailProducer
}

type ICronScheduler interface {
	Setup()
	Start()
	Stop() context.Context
}

func NewCronScheduler(
	subscriptionService services.ISubscriptionService, emailProducer *producers.EmailProducer) *CronScheduler {
	return &CronScheduler{
		Cron:                cron.New(),
		subscriptionService: subscriptionService,
		emailProducer:       emailProducer,
	}
}

func (s *CronScheduler) Setup() {
	_, err := s.Cron.AddFunc(EveryDayAt9AM, func() {
		subscribed, err := s.subscriptionService.ListSubscribed()
		if err != nil {
			log.Printf("failed to load active subscriptions: %v", err.Error())
		}

		for _, v := range subscribed {
			err = s.emailProducer.Produce(v)
			if err != nil {
				log.Printf("failed to send email to kafka: %v", err.Error())
			}
		}

	})
	if err != nil {
		log.Printf("Failed to register job: %v", err.Error())
	}
}

func (s *CronScheduler) Start() {
	s.Cron.Start()
}

func (s *CronScheduler) Stop() context.Context {
	return s.Cron.Stop()
}
