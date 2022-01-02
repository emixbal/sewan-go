package routers

import (
	"sejuta-cita/app/controllers"
	"sejuta-cita/app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Auth(app *fiber.App) {
	r := app.Group("/auth")

	r.Post("/login", controllers.UserLogin)
	r.Post("/register", controllers.UserRegister)

	r.Post("/new-password", middlewares.IsAuthenticated, controllers.NewPasswordSelf)

	r.Post("/login-n", controllers.LoginRefrehToken)
	r.Post("/refresh", controllers.RefreshToken)

}
