//go:build integration

package spire_test

import (
	"context"
	"testing"
	"time"

	"github.com/helixgitpx/platform/spire"
)

func TestNewFetcher_ConnectsToAgent(t *testing.T) {
	// Requires a SPIRE agent socket at unix:///run/spire/agent.sock
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	f, err := spire.NewFetcher(ctx, spire.Options{SocketPath: "unix:///run/spire/agent.sock"})
	if err != nil {
		t.Fatalf("NewFetcher: %v", err)
	}
	defer f.Close()
	if f.Source() == nil {
		t.Fatal("expected live Source when socket is present")
	}
}
