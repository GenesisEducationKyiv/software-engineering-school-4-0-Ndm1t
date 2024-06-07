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
	Email     string       `gorm:"primaryKey"`
	Status    Subscription `gorm:"default:subscribed"`
	CreatedAt time.Time    `gorm:"autoCreateTime"`
	DeletedAt *time.Time   `gorm:"index"`
}
