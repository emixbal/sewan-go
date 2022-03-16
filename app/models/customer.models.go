package models

import (
	"time"
)

type Customer struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email" gorm:"index:idx_email,unique"`
	Phone     string `json:"phone" gorm:"index:idx_phone,unique"`
	IsActive  bool   `json:"is_active,omitempty" gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
