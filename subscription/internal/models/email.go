package models

import (
	"time"
)

type Subscription string

const (
	Subscribed   Subscription = "subscribed"
	Unsubscribed Subscription = "unsubscribed"
)

type Email struct {
	Email     string       `gorm:"primaryKey" json:"email"`
	Status    Subscription `gorm:"default:subscribed" json:"status"`
	CreatedAt time.Time    `gorm:"autoCreateTime" json:"createdAt"`
	DeletedAt *time.Time   `gorm:"index" json:"deletedAt"`
}
