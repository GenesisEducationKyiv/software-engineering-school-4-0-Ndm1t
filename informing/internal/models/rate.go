package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Rate struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	Rate      float64   `json:"rate"`
	CreatedAt time.Time `json:"createdAt"`
}

func (r *Rate) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID, _ = uuid.NewUUID()
	return
}
