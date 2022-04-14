package models

import "time"

type Payment struct {
	ID            uint `json:"id"`
	TransactionID uint `json:"transaction_id" gorm:"not null"`
	Transaction   Transaction
	Amount        int       `json:"amount"`
	Date          time.Time `gorm:"not null" json:"date"`
}
