package database

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"informing-service/internal/models"
	"time"
)

type RateRepository struct {
	DB *gorm.DB
}

func NewRateRepository(db *gorm.DB) *RateRepository {
	return &RateRepository{
		DB: db,
	}
}

func (r *RateRepository) Create(rate models.Rate) (*models.Rate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	result := r.DB.WithContext(ctx).Create(&rate)

	if result.Error != nil {
		return nil, result.Error
	}

	return &rate, nil
}

func (r *RateRepository) GetLatest() (*models.Rate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var rate models.Rate

	result := r.DB.WithContext(ctx).Order("created_at desc").Find(&rate)

	if result.Error != nil {
		return nil, fmt.Errorf("database error: %v", result.Error)
	}

	return &rate, nil
}
