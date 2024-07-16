package database

import (
	"customer-service/internal/models"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrDatabaseInternal = errors.New("database internal error")
)

type (
	CustomerRepository struct {
		db *gorm.DB
	}
)

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (r *CustomerRepository) Create(customer models.Customer) (*models.Customer, error) {
	result := r.db.Create(&customer)
	if result.Error != nil {
		return nil, ErrDatabaseInternal
	}
	return &customer, nil
}

func (r *CustomerRepository) DeleteByTxID(txID uuid.UUID) error {
	result := r.db.Model(&models.Customer{}).Where("txId = ?", txID)
	if result.Error != nil {
		return ErrDatabaseInternal
	}
	return nil
}
