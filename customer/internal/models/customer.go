package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	ID    uuid.UUID `gorm:"primaryKey" json:"id"`
	TxID  uuid.UUID `json:"txId"`
	Email string    `json:"email"`
}

func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID, err = uuid.NewUUID()
	if err != nil {
		return err
	}
	return
}
