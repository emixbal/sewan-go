package routers

import (
	"sewan-go/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func Transaction(app *fiber.App) {
	r := app.Group("/transactions")

	r.Post("/", controllers.TransactionNew)
	r.Post("/customer", controllers.TransactionCustomerNew)
	r.Get("/:transaction_id", controllers.TransactionDetail)
	r.Get("/:transaction_id/items", controllers.TransactionShowItems)
	r.Post("/:transaction_id/add-items", controllers.AddItemToTransaction)
	r.Delete("/item/:item_id", controllers.TransactionItemDelete)
	r.Put("/item/:item_id", controllers.TransactionItemUpdateQty)
}
