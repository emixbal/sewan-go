package routers

import (
	"sejuta-cita/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func Auth(app *fiber.App) {
	r := app.Group("/auth")

	r.Post("/login", controllers.UserLogin)
	r.Post("/register", controllers.UserRegister)
	// r.Post("/make_admin", controllers.CreateBook)
	// r.Post("/remove_admin", controllers.CreateBook)
}
