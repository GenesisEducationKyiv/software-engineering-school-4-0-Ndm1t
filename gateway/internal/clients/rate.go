package clients

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const getRateEndpoint = "/rate"

type (
	RateClient struct {
		logger *zap.SugaredLogger
	}
)

func NewRateClient(logger *zap.SugaredLogger) *RateClient {
	return &RateClient{
		logger: logger,
	}
}

func (c *RateClient) FetchRate() (*int, *string, []byte, error) {
	rateUrl := viper.GetString("RATE_SERVICE_BASE_URL") + getRateEndpoint
	res, err := http.Get(rateUrl)
	if err != nil {
		return nil, nil, nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, nil, err
	}

	resContentType := res.Header.Get("Content-Type")

	return &res.StatusCode, &resContentType, body, nil
}
