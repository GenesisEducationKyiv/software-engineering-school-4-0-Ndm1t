package database

import (
	"context"
	"gorm.io/gorm"
	"subscription-service/internal/app_errors"
	"subscription-service/internal/models"
	"time"
)

type SubscriptionRepository struct {
	DB *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		DB: db,
	}
}

func (d *SubscriptionRepository) Find(email string) (*models.Email, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var subscription models.Email
	result := d.DB.WithContext(ctx).Where("email = ?", email).Find(&subscription)

	if result.Error != nil {
		return nil, apperrors.ErrDatabase
	}

	return &subscription, nil
}

func (d *SubscriptionRepository) Create(email string) (*models.Email, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DBTimeout)
	defer cancel()

	subscription := models.Email{Email: email, Status: models.Subscribed}
	result := d.DB.WithContext(ctx).Create(&subscription)

	if result.Error != nil {
		return nil, apperrors.ErrDatabase
	}

	return &subscription, nil
}

func (d *SubscriptionRepository) ListSubscribed() ([]models.Email, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DBTimeout)
	defer cancel()

	var subscriptions []models.Email
	result := d.DB.WithContext(ctx).Find(&subscriptions, "status = ?", models.Subscribed)

	if result.Error != nil {
		return nil, apperrors.ErrDatabase
	}

	return subscriptions, nil
}

func (d *SubscriptionRepository) Update(subscription models.Email) (*models.Email, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DBTimeout)
	defer cancel()

	result := d.DB.WithContext(ctx).Updates(&subscription)

	if result.Error != nil {
		return nil, apperrors.ErrDatabase
	}

	return &subscription, nil
}

func (d *SubscriptionRepository) Delete(subscription *models.Email) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	result := d.DB.WithContext(ctx).Delete(subscription)

	if result.Error != nil {
		return apperrors.ErrDatabase
	}

	return nil
}
