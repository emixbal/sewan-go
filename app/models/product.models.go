package models

import "time"

type Product struct {
	ID        uint   `json:"id"`
	IsActive  bool   `json:"is_active,omitempty" gorm:"default:true"`
	Name      string `json:"name"`
	Kode      string `json:"kode"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
