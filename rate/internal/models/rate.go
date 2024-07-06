package models

import (
	"time"
)

type Rate struct {
	Rate      float64   `json:"rate"`
	CreatedAt time.Time `json:"createdAt"`
}
