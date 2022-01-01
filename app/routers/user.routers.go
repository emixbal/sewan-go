package routers

import (
	"sejuta-cita/app/controllers"
	"sejuta-cita/app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func User(app *fiber.App) {
	r := app.Group("/user")

	r.Get("/", middlewares.ExampleMiddleware, controllers.FetchAllUsers) // contoh menggunakan middleware
	r.Post("/", controllers.CreateUser)
}
