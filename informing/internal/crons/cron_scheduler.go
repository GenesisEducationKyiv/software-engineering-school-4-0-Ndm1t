package crons

import (
	"context"
	"github.com/robfig/cron/v3"
	"informing-service/internal/services"
	"log"
)

const (
	EveryDayAt9AM = "0 9 * * *"
)

type CronScheduler struct {
	Cron             *cron.Cron
	informingService services.InformingServiceInterface
}

type ICronScheduler interface {
	Setup()
	Start()
	Stop() context.Context
}

func NewCronScheduler(
	informingService services.InformingServiceInterface) *CronScheduler {
	return &CronScheduler{
		Cron:             cron.New(),
		informingService: informingService,
	}
}

func (s *CronScheduler) Setup() {
	_, err := s.Cron.AddFunc(EveryDayAt9AM, s.informingService.SendEmails)
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
