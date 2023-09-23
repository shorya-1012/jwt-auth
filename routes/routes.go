package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shorya-1012/jwt-auth/controllers"
)

func Setup(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to jwt auth")
	})

	app.Post("/register", controllers.RegisterController)

	app.Post("/login", controllers.LoginController)

	app.Get("/token", controllers.RefreshTokenController)

	app.Delete("/logout", controllers.LogoutController)

	app.Get("/user", controllers.FetchUserController)
}
