package models

import "time"

type SubscriptionStatus string

const (
	Subscribed   SubscriptionStatus = "subscribed"
	Unsubscribed SubscriptionStatus = "unsubscribed"
)

type Subscription struct {
	Email     string             `gorm:"primaryKey" json:"email"`
	Status    SubscriptionStatus `gorm:"default:subscribed" json:"status"`
	CreatedAt time.Time          `gorm:"autoCreateTime" json:"createdAt"`
	DeletedAt *time.Time         `gorm:"index" json:"deletedAt"`
}
