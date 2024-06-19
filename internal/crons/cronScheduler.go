package crons

import (
	"github.com/robfig/cron/v3"
	"gses4_project/internal/services"
	"log"
)

type CronScheduler struct {
	Cron *cron.Cron
}

func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		Cron: cron.New(),
	}
}

func (s *CronScheduler) Setup() {
	informingService := services.NewInformingService()
	_, err := s.Cron.AddFunc("0 9 * * *", informingService.SendEmails)
	if err != nil {
		log.Printf("Failed to register job: %v", err.Error())
	}
}
