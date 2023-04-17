package main

import (
	"fmt"
	"log"

	"github.com/arvians-id/go-rabbitmq/gateway/api/user/handler"

	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	configuration := config.New()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-KEY",
		AllowMethods:     "POST, DELETE, PUT, PATCH, GET",
		AllowCredentials: true,
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to my API Todo List")
	})

	handler.NewUserHandler(app, configuration)

	port := fmt.Sprintf(":%s", configuration.Get("APP_PORT"))
	err := app.Listen(port)
	if err != nil {
		log.Fatalln("Cannot connect to server", err)
	}
}
