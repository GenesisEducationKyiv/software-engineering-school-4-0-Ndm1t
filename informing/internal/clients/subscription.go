package clients

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

const listSubscribedEndpoint = "/subscriptions"

type (
	ListSubscribedRes struct {
		Subscriptions []string `json:"subscriptions"`
	}

	SubscriptionClient struct{}
)

func NewSubscriptionClient() *SubscriptionClient {
	return &SubscriptionClient{}
}

func (c *SubscriptionClient) ListSubscribed() ([]string, error) {
	res, err := http.Get(os.Getenv("SUBSCRIPTION_SERVICE_BASE_URL") + listSubscribedEndpoint)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var subscriptionsData ListSubscribedRes
	err = json.Unmarshal(body, &subscriptionsData)
	if err != nil {
		return nil, err
	}

	return subscriptionsData.Subscriptions, nil
}
