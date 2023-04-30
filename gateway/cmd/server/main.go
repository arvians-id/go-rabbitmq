package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
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
	// Init Config
	configuration := config.New(".env.dev")

	// Init Open Telementry Tracer
	tp, err := config.NewTracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
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

	// Init Redis
	rdb, err := config.InitRedis(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// Init Server
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			return response.ReturnError(ctx, code, err)
		},
	})
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-KEY",
		AllowMethods:     "POST, DELETE, PUT, PATCH, GET",
		AllowCredentials: true,
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		tr := tp.Tracer("main-endpoint")

		_, span := tr.Start(ctx, "root")
		defer span.End()
		return c.SendString("Welcome to my API Todo List")
	})

	api.NewRoutes(app, configuration, rdb)

	port := fmt.Sprintf(":%s", configuration.Get("APP_PORT"))
	err = app.Listen(port)
	if err != nil {
		log.Fatalln("Cannot connect to server", err)
	}
}
