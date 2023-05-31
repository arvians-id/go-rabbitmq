package api

import (
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/arvians-id/go-rabbitmq/gateway/api/gql"
	"github.com/arvians-id/go-rabbitmq/gateway/api/gql/resolver"
	"github.com/arvians-id/go-rabbitmq/gateway/api/middleware"
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/auth"
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/category"
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/category_todo"
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/todo"
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/user"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/rabbitmq/amqp091-go"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"net/http"
	"os"
	"time"
)

func NewRoutes(configuration config.Config, logFile *os.File, ch *amqp091.Channel) (*fiber.App, error) {
	// Init Redis
	redisClient, err := config.InitRedis(configuration)
	if err != nil {
		panic(err)
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

	// Set CSRF
	if configuration.Get("STATE") == "production" {
		app.Use(csrf.New(csrf.Config{
			KeyLookup:      "header:X-CSRF-Token",
			CookieName:     "csrf_token",
			CookieSameSite: "Lax",
			CookieHTTPOnly: true,
			Expiration:     15 * time.Minute,
			KeyGenerator:   utils.UUID,
		}))
	}

	// Set Etag
	app.Use(etag.New())

	// Set Logging
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		Output:     logFile,
		TimeFormat: "02-Jan-2006 15:04:05",
		Done: func(c *fiber.Ctx, logString []byte) {
			fmt.Print(string(logString))
		},
	}))

	// Set CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-KEY",
		AllowMethods:     "POST, DELETE, PUT, PATCH, GET",
		AllowCredentials: true,
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to my API Todo List")
	})

	// Set Routes
	auth.NewAuthRoute(app, configuration, ch)

	apiGroup := app.Group("/api")
	userService := user.NewUserRoute(apiGroup, configuration, ch)
	categoryService := category.NewCategoryRoute(apiGroup, configuration)
	todoService := todo.NewTodoRoute(apiGroup, configuration, redisClient)
	category_todo.NewCategoryTodoRoute(apiGroup, ch)

	// Set GraphQL Playground
	app.Get("/playground", func(c *fiber.Ctx) error {
		h := playground.Handler("GraphQL", "/query")
		fasthttpadaptor.NewFastHTTPHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			h.ServeHTTP(writer, request)
		})(c.Context())

		return nil
	})

	resolvers := &resolver.Resolver{
		UserService:      userService,
		CategoryServices: categoryService,
		TodoService:      todoService,
	}

	generatedConfig := gql.Config{
		Resolvers: resolvers,
	}

	app.Use(middleware.DataLoaders(userService, todoService, categoryService))

	h := handler.NewDefaultServer(gql.NewExecutableSchema(generatedConfig))
	app.Post("/query", func(c *fiber.Ctx) error {
		fasthttpadaptor.NewFastHTTPHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			h.ServeHTTP(writer, request)
		})(c.Context())

		return nil
	})

	return app, nil
}
