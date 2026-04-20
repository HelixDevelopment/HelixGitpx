package temporal_test

import (
	"context"
	"testing"

	"github.com/helixgitpx/platform/temporal"
)

func TestNewClient_NoopWhenAddrEmpty(t *testing.T) {
	c, err := temporal.NewClient(context.Background(), temporal.Options{})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if c == nil {
		t.Fatalf("nil client")
	}
	if err := c.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}
