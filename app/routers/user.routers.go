package routers

import (
	"sejuta-cita/app/controllers"
	"sejuta-cita/app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func User(app *fiber.App) {
	IsAuthenticated := app.Group("/user",
		middlewares.IsAuthenticated,
	)

	IsAuthenticated.Get("/all",
		controllers.FetchAllUsers,
	)

	IsAdmin := IsAuthenticated.Group("/",
		middlewares.IsAdmin,
	)
	IsAdmin.Post("/new",
		controllers.UserRegister,
	)
}
