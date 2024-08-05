package clients

import (
	"bytes"
	"encoding/json"
	"gateway/internal/server/controllers"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

const (
	subscribeEndpoint   = "/subscribe"
	unsubscribeEndpoint = "/unsubscribe"
	contentType         = "application/json"
)

type (
	SubscriptionClient struct {
		logger Logger
	}
)

func NewSubscriptionClient(logger Logger) *SubscriptionClient {
	return &SubscriptionClient{logger: logger}
}

func (c *SubscriptionClient) Subscribe(req controllers.SubscribeReq) (*int, *string, []byte, error) {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		c.logger.Warnf("failed to marshal request: %v", err.Error())
		return nil, nil, nil, err
	}
	res, err := http.Post(
		viper.GetString("SUBSCRIPTION_SERVICE_BASE_URL")+subscribeEndpoint,
		contentType,
		bytes.NewBuffer(reqJSON))
	if err != nil {
		c.logger.Warnf("failed to send post request: %v", err.Error())
		return nil, nil, nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Warnf("failed to read response body: %v", err.Error())
		return nil, nil, nil, err
	}

	resContentType := res.Header.Get("Content-Type")

	return &res.StatusCode, &resContentType, body, nil
}

func (c *SubscriptionClient) Unsubscribe(req controllers.SubscribeReq) (*int, *string, []byte, error) {

	reqJSON, err := json.Marshal(req)
	if err != nil {
		c.logger.Warnf("failed to marshal request: %v", err.Error())
		return nil, nil, nil, err
	}
	res, err := http.Post(
		viper.GetString("SUBSCRIPTION_SERVICE_BASE_URL")+unsubscribeEndpoint,
		contentType,
		bytes.NewBuffer(reqJSON))
	if err != nil {
		c.logger.Warnf("failed to send post request: %v", err.Error())
		return nil, nil, nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Warnf("failed to read response body: %v", err.Error())
		return nil, nil, nil, err
	}

	resContentType := res.Header.Get("Content-Type")

	return &res.StatusCode, &resContentType, body, nil
}
