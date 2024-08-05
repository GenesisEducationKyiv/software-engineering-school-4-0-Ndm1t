package models

import "github.com/google/uuid"

type Subscription struct {
	TxID  uuid.UUID `json:"txId"`
	Email string    `json:"email"`
}
