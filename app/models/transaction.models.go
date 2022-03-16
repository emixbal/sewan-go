package models

import "time"

type Transaction struct {
	ID         uint `json:"id"`
	IsActive   bool `json:"is_active,omitempty" gorm:"default:true"`
	CustomerID int  `json:"customer_id"`
	Customer   Customer
	StartDate  time.Time `gorm:"not null" json:"start_date"`
	EndDate    time.Time `gorm:"not null" json:"end_date"`
	CreatedAt  time.Time `json:"created_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP" faker:"-"`
	UpdatedAt  time.Time `json:"update_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" faker:"-"`
}
