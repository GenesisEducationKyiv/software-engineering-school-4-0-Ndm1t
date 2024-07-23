package providers

import (
	"github.com/VictoriaMetrics/metrics"
	"go.uber.org/zap"
	"rate-service/internal/services"
)

const (
	rateGaugeMetric = "rate_gauge"
)

type (
	LoggingProvider struct {
		name         string
		rateProvider services.IRateAPIProvider
		logger       *zap.SugaredLogger
	}
)

func NewLoggingProvider(name string, rateProvider services.IRateAPIProvider, logger *zap.SugaredLogger) *LoggingProvider {
	return &LoggingProvider{
		name:         name,
		rateProvider: rateProvider,
		logger:       logger,
	}
}
func (l *LoggingProvider) FetchRate() (*float64, error) {
	rate, err := l.rateProvider.FetchRate()
	if err != nil {
		l.logger.Warnf("%v provider returned error: %v", l.name, err)
		return nil, err
	}
	l.logger.Infof("%v provider returned value: %v", l.name, *rate)
	metrics.GetOrCreateGauge(rateGaugeMetric, func() float64 {
		return *rate
	})
	return rate, err
}
