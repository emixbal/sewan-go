package routers

import (
	"sejuta-cita/app/controllers"
	"sejuta-cita/app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Book(app *fiber.App) {
	r := app.Group("/books")

	r.Get("/", middlewares.ExampleMiddleware, controllers.FetchAllBooks) // contoh menggunakan middleware
	r.Post("/", controllers.CreateBook)
}
