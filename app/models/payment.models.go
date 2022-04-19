package models

import "time"

type Payment struct {
	ID            uint `json:"id"`
	TransactionID uint `json:"transaction_id" gorm:"not null"`
	Transaction   Transaction
	Nominal       int       `json:"nominal"`
	Date          time.Time `gorm:"not null" json:"date"`
}
