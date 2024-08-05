package models

import "github.com/google/uuid"

type Customer struct {
	TxID  uuid.UUID `json:"txId"`
	Email string    `json:"email"`
}
