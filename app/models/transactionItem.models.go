package models

import "time"

type TransactionItem struct {
	ID            uint `json:"id"`
	TransactionID int  `json:"transaction_id" gorm:"not null"`
	Transaction   Transaction
	ProductID     int `gorm:"not null" json:"product_id"`
	Product       Product
	Qty           int       `gorm:"not null" json:"qty"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP" faker:"-"`
	UpdatedAt     time.Time `json:"update_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" faker:"-"`
}
