package crons

import (
	"context"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"rate-service/internal/services"
)

const (
	EveryDayAt1Am = "0 1 * * *"
)

type CronScheduler struct {
	Cron        *cron.Cron
	rateService services.IRateService
	logger      *zap.SugaredLogger
}

type ICronScheduler interface {
	Setup()
	Start()
	Stop() context.Context
}

func NewCronScheduler(
	rateService services.IRateService,
	logger *zap.SugaredLogger) *CronScheduler {
	return &CronScheduler{
		Cron:        cron.New(),
		rateService: rateService,
		logger:      logger,
	}
}

func (s *CronScheduler) Setup() {
	_, err := s.Cron.AddFunc(EveryDayAt1Am, func() {
		_, err := s.rateService.Get()
		if err != nil {
			s.logger.Warn(err)
		}
	})
	if err != nil {
		s.logger.Warnf("Failed to register job: %v", err.Error())
	}
}

func (s *CronScheduler) Start() {
	s.Cron.Start()
}

func (s *CronScheduler) Stop() context.Context {
	return s.Cron.Stop()
}
