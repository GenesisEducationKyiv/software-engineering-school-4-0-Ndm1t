package database

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"informing-service/internal/models"
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

func (d *SubscriptionRepository) Create(email string) (*models.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	subscription := models.Subscription{Email: email, Status: models.Subscribed}
	result := d.DB.WithContext(ctx).Create(&subscription)

	if result.Error != nil {
		return nil, result.Error
	}

	return &subscription, nil
}

func (d *SubscriptionRepository) ListSubscribed() ([]models.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var subscriptions []models.Subscription
	result := d.DB.WithContext(ctx).Find(&subscriptions, "status = ?", models.Subscribed)

	if result.Error != nil {
		return nil, fmt.Errorf("databese error: %v", result.Error.Error())
	}

	return subscriptions, nil
}

func (d *SubscriptionRepository) Update(subscription models.Subscription) (*models.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	result := d.DB.WithContext(ctx).Updates(&subscription)

	if result.Error != nil {
		return nil, fmt.Errorf("databese error: %v", result.Error.Error())
	}

	return &subscription, nil
}
