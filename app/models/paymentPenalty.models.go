package models

import "time"

type PaymentPenalty struct {
	ID            uint `json:"id"`
	TransactionID uint `json:"transaction_id" gorm:"not null"`
	Transaction   Transaction
	Nominal       int       `json:"nominal"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP" faker:"-"`
}
