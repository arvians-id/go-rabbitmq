package main

import (
	"context"
	"fmt"
	"github.com/arvians-id/go-rabbitmq/gateway/api"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log"
	"os"
	"time"
)

func main() {
	// Init Config
	configuration := config.New()

	// Set Context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Init Open Telementry Tracer
	tp, err := config.NewTracerProvider(configuration)
	if err != nil {
		log.Fatalln("There is something wrong with the tracer provider", err)
	}

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

	// Init Log File
	file, err := os.OpenFile("./logs/main.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("There is something wrong with the log file", err)
	}
	defer file.Close()

	// Start Server
	app, err := api.NewRoutes(configuration, file)
	if err != nil {
		log.Fatalln("There is something wrong with the server", err)
	}

	port := fmt.Sprintf(":%s", configuration.Get("APP_PORT"))
	err = app.Listen(port)
	if err != nil {
		log.Fatalln("Cannot connect to server", err)
	}
}
