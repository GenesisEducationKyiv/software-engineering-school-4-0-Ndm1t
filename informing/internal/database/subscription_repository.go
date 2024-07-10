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

func (r *SubscriptionRepository) Create(email string) (*models.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DBTimeout)
	defer cancel()

	subscription := models.Subscription{Email: email, Status: models.Subscribed}
	result := r.DB.WithContext(ctx).Create(&subscription)

	if result.Error != nil {
		return nil, result.Error
	}

	return &subscription, nil
}

func (r *SubscriptionRepository) ListSubscribed(limit int, startEmail string) ([]models.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var subscriptions []models.Subscription
	if startEmail == "" {
		result := r.DB.WithContext(ctx).
			Order("email").Limit(limit).
			Find(&subscriptions, "status = ?", models.Subscribed)
		return subscriptions, result.Error
	}

	result := r.DB.WithContext(ctx).Where("email > ?", startEmail).
		Order("email").
		Limit(limit).
		Find(&subscriptions, "status = ?", models.Subscribed)
	return subscriptions, result.Error
}

func (r *SubscriptionRepository) Update(subscription models.Subscription) (*models.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DBTimeout)
	defer cancel()

	result := r.DB.WithContext(ctx).Updates(&subscription)

	if result.Error != nil {
		return nil, fmt.Errorf("databese error: %v", result.Error.Error())
	}

	return &subscription, nil
}

func (r *SubscriptionRepository) Delete(subscription models.Subscription) error {
	ctx, cancel := context.WithTimeout(context.Background(), DBTimeout)
	defer cancel()

	result := r.DB.WithContext(ctx).Unscoped().Delete(&subscription)

	if result.Error != nil {
		return fmt.Errorf("failed to delete subscription: %v", result.Error)
	}

	return nil
}

func (r *SubscriptionRepository) Find(email string) (*models.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var subscription models.Subscription
	result := r.DB.WithContext(ctx).Where("email = ?", email).Find(&subscription)

	if result.Error != nil {
		return nil, fmt.Errorf("databese error: %v", result.Error.Error())
	}

	return &subscription, nil
}
