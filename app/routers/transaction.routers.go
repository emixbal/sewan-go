package routers

import (
	"sewan-go/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func Transaction(app *fiber.App) {
	r := app.Group("/transactions")

	r.Post("/", controllers.NewTransaction)
}
