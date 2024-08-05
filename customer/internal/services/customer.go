package services

import (
	"context"
	"customer-service/internal/models"
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/google/uuid"
)

const (
	customerCreateSuccess = "CustomerCreated"
	customerCreatedFailed = "CustomerFailed"
	messagePublishFail    = "message_publish_fail"
	messagePublishSuccess = "message_publish_success"
)

type (
	CustomerRepositoryInterface interface {
		Create(customer models.Customer) (*models.Customer, error)
		DeleteByTxID(txId uuid.UUID) error
	}

	CustomerPublisherInterface interface {
		Publish(eventType string, customer models.Customer, ctx context.Context) error
	}

	CustomerServiceInterface interface {
		Create(txID uuid.UUID, email string)
	}

	CustomerService struct {
		publisher  CustomerPublisherInterface
		repository CustomerRepositoryInterface
	}
)

func NewCustomerService(publisher CustomerPublisherInterface, repository CustomerRepositoryInterface) *CustomerService {
	return &CustomerService{
		publisher:  publisher,
		repository: repository,
	}
}

func (s *CustomerService) Create(txID uuid.UUID, email string) {
	customer, err := s.repository.Create(models.Customer{
		TxID:  txID,
		Email: email,
	})
	if err != nil {
		err = s.publisher.Publish(customerCreatedFailed, models.Customer{}, context.Background())
		if err != nil {
			metrics.GetOrCreateCounter(fmt.Sprintf(`%v{type=%q}`, messagePublishFail, customerCreatedFailed))
		}
		metrics.GetOrCreateCounter(fmt.Sprintf(`%v{type=%q}`, messagePublishSuccess, customerCreatedFailed))
		return
	}

	err = s.publisher.Publish(customerCreateSuccess, *customer, context.Background())
	if err != nil {
		metrics.GetOrCreateCounter(fmt.Sprintf(`%v{type=%q}`, messagePublishFail, customerCreatedFailed))
	}
	metrics.GetOrCreateCounter(fmt.Sprintf(`%v{type=%q}`, messagePublishSuccess, customerCreatedFailed))
}
