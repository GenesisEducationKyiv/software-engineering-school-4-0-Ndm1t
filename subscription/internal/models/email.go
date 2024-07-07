package models

import (
	"gorm.io/gorm"
	"time"
)

type Subscription string

const (
	Subscribed   Subscription = "subscribed"
	Unsubscribed Subscription = "unsubscribed"
)

type Email struct {
	Email     string         `gorm:"primaryKey" json:"email"`
	Status    Subscription   `gorm:"default:subscribed" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (e *Email) BeforeDelete(tx *gorm.DB) (err error) {
	e.Status = Unsubscribed
	return
}
