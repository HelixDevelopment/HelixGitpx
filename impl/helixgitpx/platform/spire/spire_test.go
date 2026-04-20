package spire_test

import (
	"context"
	"testing"

	"github.com/helixgitpx/platform/spire"
)

func TestNewFetcher_NoopWhenSocketAbsent(t *testing.T) {
	f, err := spire.NewFetcher(context.Background(), spire.Options{
		SocketPath: "unix:///tmp/definitely-not-here.sock",
	})
	if err != nil {
		t.Fatalf("NewFetcher: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}
