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
	r.Get("/", controllers.TransactionList)
	r.Get("/:transaction_id/items", controllers.TransactionShowItems)
	r.Get("/:transaction_id/payment", controllers.TransactionDetailAndPayments)
	r.Post("/:transaction_id/add-items", controllers.AddItemToTransaction)
	r.Post("/:transaction_id/payment", controllers.TransactionAddPayment)
	r.Delete("/item/:item_id", controllers.TransactionItemDelete)
	r.Put("/item/:item_id", controllers.TransactionItemUpdateQty)
}
