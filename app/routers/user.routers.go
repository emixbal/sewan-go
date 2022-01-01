package routers

import (
	"sejuta-cita/app/controllers"
	"sejuta-cita/app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func User(app *fiber.App) {
	// need login access (regular users)
	IsAuthenticated := app.Group("/user",
		middlewares.IsAuthenticated,
	)

	IsAuthenticated.Get("/all",
		controllers.FetchAllUsers,
	)

	IsAuthenticated.Get("/:user_id",
		controllers.ShowUserDetail,
	)

	// need admin access
	IsAdmin := IsAuthenticated.Group("/",
		middlewares.IsAdmin,
	)
	IsAdmin.Post("/new",
		controllers.UserRegister,
	)
	IsAdmin.Delete("/:id",
		controllers.UserSoftDelete,
	)
	IsAdmin.Delete("/:id/force",
		controllers.UserHardDelete,
	)
	IsAdmin.Put("/:id",
		controllers.UserUpdate,
	)
}
