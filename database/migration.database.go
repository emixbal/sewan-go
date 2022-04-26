package database

import (
	"sewan-go/app/models"

	"gorm.io/gorm"
)

func InitMigration(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Customer{},
		&models.Transaction{},
		&models.TransactionItem{},
		&models.Payment{},
		&models.StatusTransaction{},
		&models.Demage{},
		&models.PaymentPenalty{},
	)
	InitSeeding(db)
}
func InitSeeding(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		UserSeeder(db)
	}
	db.Model(&models.StatusTransaction{}).Count(&count)
	if count == 0 {
		StatusTransactionSeeder(db)
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

func StatusTransactionSeeder(db *gorm.DB) {
	status_transaction := []models.StatusTransaction{
		{
			Name: "Default",
		},
		{
			Name: "Dikirim",
		},
		{
			Name: "Ada kerusakan/kehilangan",
		},
		{
			Name: "Selesai & Dikembalikan",
		},
	}
	db.Create(&status_transaction)
}
