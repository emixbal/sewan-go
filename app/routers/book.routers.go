package routers

import (
	"sejuta-cita/app/controllers"
	"sejuta-cita/app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Book(app *fiber.App) {
	user := app.Group("/books")

	user.Get("/", middlewares.ExampleMiddleware, controllers.FetchAllBooks) // contoh menggunakan middleware
	user.Post("/", controllers.CreateBook)
}
