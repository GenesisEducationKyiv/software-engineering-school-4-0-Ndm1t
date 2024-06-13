package database

import (
	"gses4_project/internal/apperrors"
	"gses4_project/internal/models"
	"net/http"
)

type SubscriptionDao struct{}

func NewSubcscriptionDao() *SubscriptionDao {
	return &SubscriptionDao{}
}

func (d *SubscriptionDao) Find(email string) (*models.Email, error) {
	var subscription models.Email
	result := DB.Where("email = ?", email).Find(&subscription)
	if result.Error != nil {
		return nil, apperrors.NewHttpError("Database error", http.StatusInternalServerError)
	}

	return &subscription, nil

}

func (d *SubscriptionDao) Create(email string) (*models.Email, error) {
	subscription := models.Email{Email: email, Status: models.Subscribed}
	result := DB.Create(&subscription)

	if result.Error != nil {
		return nil, apperrors.NewHttpError("Database error", http.StatusInternalServerError)
	}

	return &subscription, nil
}

func (d *SubscriptionDao) ListSubscribed() ([]models.Email, error) {
	var subscriptions []models.Email
	result := DB.Find(&subscriptions, "status = ?", models.Subscribed)
	if result.Error != nil {
		return nil, apperrors.NewHttpError("Database error", http.StatusInternalServerError)
	}

	return subscriptions, nil
}
func (d *SubscriptionDao) Update(subscription models.Email) (*models.Email, error) {
	result := DB.Updates(&subscription)
	if result.Error != nil {
		return nil, apperrors.NewHttpError("Database error", http.StatusInternalServerError)
	}

	return &subscription, nil
}
