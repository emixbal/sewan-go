package routers

import (
	"sewan-go/app/controllers"
	"sewan-go/app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Auth(app *fiber.App) {
	r := app.Group("/auth")

	r.Post("/login", controllers.LoginRefrehToken)
	r.Post("/register", controllers.UserRegister)

	r.Post("/new-password", middlewares.IsAuthenticated, controllers.NewPasswordSelf)

	r.Post("/login-n", controllers.UserLogin)
	r.Post("/refresh", controllers.RefreshToken)
}
