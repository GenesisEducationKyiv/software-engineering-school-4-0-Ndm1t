package database

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), DBTimeout)
	defer cancel()
	result := r.db.WithContext(ctx).Create(&customer)
	if result.Error != nil {
		return nil, ErrDatabaseInternal
	}
	return &customer, nil
}

func (r *CustomerRepository) DeleteByTxID(txID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), DBTimeout)
	defer cancel()
	result := r.db.WithContext(ctx).Model(&models.Customer{}).Where("txId = ?", txID)
	if result.Error != nil {
		return ErrDatabaseInternal
	}
	return nil
}
