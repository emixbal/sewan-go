package models

import "time"

type TransactionItem struct {
	ID           uint `json:"id"`
	TrasactionID int  `json:"transaction_id"`
	Transaction  Transaction
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
