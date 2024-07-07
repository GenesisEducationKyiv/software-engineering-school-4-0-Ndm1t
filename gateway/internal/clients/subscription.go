package clients

import (
	"bytes"
	"encoding/json"
	"gateway/internal/server/controllers"
	"io"
	"net/http"
	"os"
)

const (
	subscribeEndpoint   = "/subscribe"
	unsubscribeEndpoint = "/unsubscribe"
	contentType         = "application/json"
)

type (
	SubscriptionClient struct{}
)

func NewSubscriptionClient() *SubscriptionClient {
	return &SubscriptionClient{}
}

func (c *SubscriptionClient) Subscribe(req controllers.SubscribeReq) (*int, *string, []byte, error) {
	url := os.Getenv("SUBSCRIPTION_SERVICE_BASE_URL") + subscribeEndpoint
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return nil, nil, nil, err
	}
	res, err := http.Post(url, contentType, bytes.NewBuffer(reqJSON))
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

func (c *SubscriptionClient) Unsubscribe(req controllers.SubscribeReq) (*int, *string, []byte, error) {
	url := os.Getenv("SUBSCRIPTION_SERVICE_BASE_URL") + unsubscribeEndpoint
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return nil, nil, nil, err
	}
	res, err := http.Post(url, contentType, bytes.NewBuffer(reqJSON))
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
