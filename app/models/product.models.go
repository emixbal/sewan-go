package models

import "time"

type Product struct {
	ID        uint      `json:"id"`
	IsActive  bool      `json:"is_active,omitempty" gorm:"default:true"`
	Name      string    `json:"name"`
	Kode      string    `json:"kode"`
	CreatedAt time.Time `json:"created_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP" faker:"-"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP" faker:"-"`
}
