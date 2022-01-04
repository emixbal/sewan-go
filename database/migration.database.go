package database

import (
	"sejuta-cita/app/models"

	"gorm.io/gorm"
)

func InitMigration(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
	InitSeeding(db)
}
func InitSeeding(db *gorm.DB) {
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		UserSeeder(db)
	}
}
func UserSeeder(db *gorm.DB) {
	user := []models.User{
		{
			Name:     "muhammad iqbal",
			Email:    "emixbal@gmail.com",
			Password: "$2a$10$xO0eiq3.64vo1gR1cKkEE.hwn0OvafrzVI0HhsZWeb9UuXsl7bZrq", //aaaaaaaa
			IsAdmin:  true,
		},
	}
	db.Create(&user)
}
