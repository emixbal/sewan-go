package routers

import (
	"sewan-go/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func Product(app *fiber.App) {
	r := app.Group("/products")

	r.Get("/", controllers.FetchAllproducts)
	r.Get("/:product_id", controllers.ShowProductDetail)
	r.Post("/", controllers.CreateANewProduct)
	r.Put("/:product_id", controllers.UpdateProduct)
	r.Delete("/:product_id", controllers.ProductSoftDelete)
}
