// Package telemetry bootstraps OpenTelemetry SDK. When OTEL_EXPORTER_OTLP_ENDPOINT
// is unset (or Options.OTLPEndpoint is empty), the SDK installs a no-op tracer
// and meter, so calling Start on developer machines without a collector is safe.
package telemetry

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Options configures Start.
type Options struct {
	Service      string
	Version      string
	Environment  string
	OTLPEndpoint string // overrides OTEL_EXPORTER_OTLP_ENDPOINT when set
}

// ShutdownFunc flushes and closes all providers.
type ShutdownFunc func(context.Context) error

// Start installs global TracerProvider. Returns a no-op shutdown when no endpoint.
func Start(ctx context.Context, opts Options) (ShutdownFunc, error) {
	endpoint := opts.OTLPEndpoint
	if endpoint == "" {
		endpoint = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	}
	if endpoint == "" {
		// No collector → no-op. Tracing calls become cheap no-ops through the global.
		return func(context.Context) error { return nil }, nil
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(opts.Service),
			semconv.ServiceVersion(opts.Version),
			semconv.DeploymentEnvironment(opts.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpointURL(endpoint),
		otlptracegrpc.WithInsecure(),
	))
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(5*time.Second)),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return func(shutdownCtx context.Context) error {
		return tp.Shutdown(shutdownCtx)
	}, nil
}
