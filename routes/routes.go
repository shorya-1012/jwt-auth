package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shorya-1012/jwt-auth/controllers"
)

func Setup(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to jwt auth")
	})

	app.Post("/api/register", controllers.RegisterController)

	app.Post("/api/login", controllers.LoginController)

	app.Get("/api/token", controllers.RefreshTokenController)

	app.Delete("/api/logout", controllers.LogoutController)

	app.Get("/api/user", controllers.FetchUserController)
}
