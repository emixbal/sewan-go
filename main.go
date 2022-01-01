package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"sejuta-cita/app/routers"
	"sejuta-cita/config"
	"sejuta-cita/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	config.InitDB()
	db := config.DB
	database.InitMigration(db)

	routers.Init(app)
	app.Listen(":" + os.Getenv("APP_PORT"))
}
