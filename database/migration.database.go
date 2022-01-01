package database

import (
	"sejuta-cita/app/models"

	"gorm.io/gorm"
)

func InitMigration(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Book{})
}
