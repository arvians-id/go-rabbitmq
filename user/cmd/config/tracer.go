package config

import (
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func NewTracerProvider(configuration Config) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	url := fmt.Sprintf("http://%s:%s/api/traces", configuration.Get("JAEGER_HOST"), configuration.Get("JAEGER_PORT"))
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(ServiceTrace),
			attribute.String("environment", EnvironmentTrace),
		)),
		tracesdk.WithSampler(trace.ParentBased(trace.AlwaysSample())),
	)
	return tp, nil
}

const (
	ServiceTrace     = "user-tracer"
	EnvironmentTrace = "development"
)
