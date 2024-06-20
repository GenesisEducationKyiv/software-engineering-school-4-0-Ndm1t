package crons

import (
	"github.com/robfig/cron/v3"
	"gses4_project/internal/container"
	"gses4_project/internal/services"
	"log"
)

const (
	EveryDayAt9AM = "0 9 * * *"
)

type CronScheduler struct {
	Cron      *cron.Cron
	container container.IContainer
}

func NewCronScheduler(container container.IContainer) *CronScheduler {
	return &CronScheduler{
		Cron:      cron.New(),
		container: container,
	}
}

func (s *CronScheduler) Setup() {
	informingService := services.NewInformingService(s.container)
	_, err := s.Cron.AddFunc(EveryDayAt9AM, informingService.SendEmails)
	if err != nil {
		log.Printf("Failed to register job: %v", err.Error())
	}
}
