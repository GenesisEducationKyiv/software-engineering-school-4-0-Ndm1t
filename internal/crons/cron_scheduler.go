package crons

import (
	"context"
	"github.com/robfig/cron/v3"
	"gses4_project/internal/container"
	"gses4_project/internal/services"
	"log"
)

const (
	EveryDayAt9AM = "0 9 * * *"
)

type CronScheduler struct {
	Cron             *cron.Cron
	container        container.IContainer
	informingService services.IInformingService
}

type ICronScheduler interface {
	Setup()
	Start()
	Stop() context.Context
}

func NewCronScheduler(container container.IContainer,
	informingService services.IInformingService) *CronScheduler {
	return &CronScheduler{
		Cron:             cron.New(),
		container:        container,
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
