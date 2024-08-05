package chain

import (
	"rate-service/internal/services"
)

type (
	Chain interface {
		services.IRateAPIProvider
		SetNext(chain Chain)
	}

	BaseChain struct {
		rateProvider services.IRateAPIProvider
		next         Chain
	}
)

func NewBaseChain(rateProvider services.IRateAPIProvider) *BaseChain {
	return &BaseChain{
		rateProvider: rateProvider,
	}
}

func (c *BaseChain) SetNext(chain Chain) {
	c.next = chain
}

func (c *BaseChain) FetchRate() (*float64, error) {
	rate, err := c.rateProvider.FetchRate()
	if err != nil {
		if c.next == nil {
			return nil, err
		}
		return c.next.FetchRate()
	}
	return rate, nil
}
