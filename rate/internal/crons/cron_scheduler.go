package crons

import (
	"context"
	"github.com/robfig/cron/v3"
	"log"
	"rate-service/internal/services"
)

const (
	EveryDayAt1Am = "0 1 * * *"
)

type CronScheduler struct {
	Cron        *cron.Cron
	rateService services.IRateService
}

type ICronScheduler interface {
	Setup()
	Start()
	Stop() context.Context
}

func NewCronScheduler(
	rateService services.IRateService) *CronScheduler {
	return &CronScheduler{
		Cron:        cron.New(),
		rateService: rateService,
	}
}

func (s *CronScheduler) Setup() {
	_, err := s.Cron.AddFunc(EveryDayAt1Am, func() {
		_, err := s.rateService.Get()
		if err != nil {
			log.Print(err)
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
