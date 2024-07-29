package clients

import (
	"github.com/spf13/viper"
	"io"
	"net/http"
)

const getRateEndpoint = "/rate"

type (
	Logger interface {
		Warnf(template string, arguments ...interface{})
	}

	RateClient struct {
		logger Logger
	}
)

func NewRateClient(logger Logger) *RateClient {
	return &RateClient{
		logger: logger,
	}
}

func (c *RateClient) FetchRate() (*int, *string, []byte, error) {
	res, err := http.Get(viper.GetString("RATE_SERVICE_BASE_URL") + getRateEndpoint)
	if err != nil {
		c.logger.Warnf(`failed to fetch rate: %v`, err.Error())
		return nil, nil, nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Warnf(`failed to read response body: %v`, err.Error())
		return nil, nil, nil, err
	}

	resContentType := res.Header.Get("Content-Type")

	return &res.StatusCode, &resContentType, body, nil
}
