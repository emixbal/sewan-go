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
	r.Put("/:transaction_id/send", controllers.TransactionChangeSendStatus)
	r.Put("/:transaction_id/return/nok", controllers.TransactionChangeReturnStatusNOK)
	r.Put("/:transaction_id/return/ok", controllers.TransactionChangeReturnStatusOK)
	r.Post("/:transaction_id/add-items", controllers.AddItemToTransaction)
	r.Post("/:transaction_id/payment", controllers.TransactionAddPayment)
	r.Delete("/item/:item_id", controllers.TransactionItemDelete)
	r.Put("/item/:item_id", controllers.TransactionItemUpdateQty)

	r.Post("/:transaction_id/add-demage", controllers.TransactionAddDemage)
	r.Get("/:transaction_id/list-demage", controllers.TransactionListDemage)
}
