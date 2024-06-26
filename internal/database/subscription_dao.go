package database

import (
	"gorm.io/gorm"
	"gses4_project/internal/apperrors"
	"gses4_project/internal/models"
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
	var subscription models.Email
	result := d.DB.Where("email = ?", email).Find(&subscription)
	if result.Error != nil {
		return nil, apperrors.ErrDatabase
	}

	return &subscription, nil

}

func (d *SubscriptionRepository) Create(email string) (*models.Email, error) {
	subscription := models.Email{Email: email, Status: models.Subscribed}
	result := d.DB.Create(&subscription)

	if result.Error != nil {
		return nil, apperrors.ErrDatabase
	}

	return &subscription, nil
}

func (d *SubscriptionRepository) ListSubscribed() ([]models.Email, error) {
	var subscriptions []models.Email
	result := d.DB.Find(&subscriptions, "status = ?", models.Subscribed)
	if result.Error != nil {
		return nil, apperrors.ErrDatabase
	}

	return subscriptions, nil
}
func (d *SubscriptionRepository) Update(subscription models.Email) (*models.Email, error) {
	result := d.DB.Updates(&subscription)
	if result.Error != nil {
		return nil, apperrors.ErrDatabase
	}

	return &subscription, nil
}
