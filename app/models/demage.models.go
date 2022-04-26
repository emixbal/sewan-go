package models

import "time"

type Demage struct {
	ID            uint `json:"id"`
	TransactionID int  `json:"transaction_id" gorm:"not null"`
	Transaction   Transaction
	ProductID     int `json:"product_id" gorm:"not null"`
	Product       Product
	Qty           int       `gorm:"not null" json:"qty"`
	IsActive      bool      `json:"is_active,omitempty" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP" faker:"-"`
	UpdatedAt     time.Time `json:"update_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" faker:"-"`
}
