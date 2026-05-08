package telemetry_test

import (
	"context"
	"os"
	"testing"

	"github.com/helixgitpx/platform/telemetry"
)

func TestStart_NoEndpoint_ReturnsNoop(t *testing.T) {
	ctx := context.Background()
	shutdown, err := telemetry.Start(ctx, telemetry.Options{Service: "hello"})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	if err := shutdown(ctx); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
}

func TestStart_WithOTLPEndpoint_ConnectsAndShutsDown(t *testing.T) {
	addr := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if addr == "" {
		addr = os.Getenv("TEST_OTLP_ENDPOINT")
	}
	if addr == "" {
		t.Skip("SKIP-OK: #OTLP — set OTEL_EXPORTER_OTLP_ENDPOINT to run")
	}

	ctx := context.Background()
	shutdown, err := telemetry.Start(ctx, telemetry.Options{
		Service:      "test-service",
		Version:      "test",
		OTLPEndpoint: addr,
	})
	if err != nil {
		t.Fatalf("Start with OTLP: %v", err)
	}
	if shutdown == nil {
		t.Fatal("shutdown func must not be nil")
	}
	if err := shutdown(ctx); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
}
