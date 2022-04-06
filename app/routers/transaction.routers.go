package routers

import (
	"sewan-go/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func Transaction(app *fiber.App) {
	r := app.Group("/transactions")

	r.Post("/", controllers.TransactionNew)
	r.Get("/:transaction_id", controllers.TransactionDetail)
	r.Post("/:transaction_id/add-items", controllers.AddItemToTransaction)

	r.Delete("/:item_id", controllers.TransactionItemDelete)
}
