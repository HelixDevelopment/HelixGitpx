package telemetry_test

import (
	"context"
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
