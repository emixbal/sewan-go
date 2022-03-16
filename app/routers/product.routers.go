package routers

import (
	"sewan-go/app/controllers"
	"sewan-go/app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Product(app *fiber.App) {
	r := app.Group("/products")

	r.Get("/", middlewares.ExampleMiddleware, controllers.FetchAllproducts)
	r.Post("/", middlewares.ExampleMiddleware, controllers.CreateANewProduct)
}
