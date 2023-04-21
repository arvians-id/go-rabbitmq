package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/arvians-id/go-rabbitmq/gateway/api"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {
	tp, err := config.NewTracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func(ctx context.Context) {
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		err := tp.Shutdown(ctx)
		if err != nil {
			log.Fatalln(err)
		}
	}(ctx)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	configuration := config.New(".env.dev")
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-KEY",
		AllowMethods:     "POST, DELETE, PUT, PATCH, GET",
		AllowCredentials: true,
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		tr := tp.Tracer("component-main")

		_, span := tr.Start(ctx, "foo")
		defer span.End()
		return c.SendString("Welcome to my API Todo List")
	})

	api.NewRoutes(app, configuration)

	port := fmt.Sprintf(":%s", configuration.Get("APP_PORT"))
	err = app.Listen(port)
	if err != nil {
		log.Fatalln("Cannot connect to server", err)
	}
}
