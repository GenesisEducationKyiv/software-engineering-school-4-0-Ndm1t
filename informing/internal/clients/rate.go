package clients

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

const getRateEndpoint = "/rate"

type (
	RateClient struct{}
)

func NewRateClient() *RateClient {
	return &RateClient{}
}

func (c *RateClient) FetchRate() (*float64, error) {
	rateUrl := os.Getenv("RATE_SERVICE_BASE_URL") + getRateEndpoint
	res, err := http.Get(rateUrl)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var rate float64
	err = json.Unmarshal(body, &rate)
	if err != nil {
		return nil, err
	}

	return &rate, nil
}
