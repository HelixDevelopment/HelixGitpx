package app

import (
	"context"
	"testing"
	"time"

	"github.com/helixgitpx/platform/log"
)

func TestRun_StartsAndShutsDown(t *testing.T) {
	t.Setenv("GIT_INGRESS_HTTP_ADDR", "127.0.0.1:0")

	ctx, cancel := context.WithCancel(context.Background())
	lg := log.Default()

	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, lg)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not shut down within 5s")
	}
}

func TestEnvOrDefault(t *testing.T) {
	if got := envOrDefault("NOT_SET_XYZ", ":8080"); got != ":8080" {
		t.Fatalf("want :8080 got %q", got)
	}
	t.Setenv("TEST_GIT_INGRESS_KEY", "  :9999  ")
	if got := envOrDefault("TEST_GIT_INGRESS_KEY", ":8080"); got != ":9999" {
		t.Fatalf("want :9999 got %q", got)
	}
}
