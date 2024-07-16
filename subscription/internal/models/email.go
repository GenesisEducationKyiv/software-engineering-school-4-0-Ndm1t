package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Subscription string

const (
	Subscribed   Subscription = "subscribed"
	Unsubscribed Subscription = "unsubscribed"
)

type State string

const (
	Pending  State = "pending"
	Verified State = "verified"
)

type Email struct {
	Email     string         `gorm:"primaryKey" json:"email"`
	Status    Subscription   `gorm:"default:subscribed" json:"status"`
	State     State          `gorm:"default:pending"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	TxID      uuid.UUID      `json:"txId"`
}

func (e *Email) BeforeDelete(tx *gorm.DB) (err error) {
	e.Status = Unsubscribed
	return
}

func (e *Email) BeforeCreate(tx *gorm.DB) (err error) {
	e.TxID, err = uuid.NewUUID()
	return
}
