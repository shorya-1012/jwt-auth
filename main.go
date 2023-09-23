package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/shorya-1012/jwt-auth/controllers"
	"github.com/shorya-1012/jwt-auth/routes"
)

func main(){
    controllers.ConnectToDB()

    app := fiber.New()

    app.Use(cors.New(cors.Config{
        AllowCredentials: true,
    }))
    routes.Setup(app)

    app.Listen(":8000")
}
