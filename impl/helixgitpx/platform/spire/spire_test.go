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
	if f.Source() != nil {
		t.Error("Source() should be nil for noop fetcher")
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}

func TestNewFetcher_EmptySocketPath(t *testing.T) {
	f, err := spire.NewFetcher(context.Background(), spire.Options{
		SocketPath: "",
	})
	if err != nil {
		t.Fatalf("NewFetcher with empty socket: %v", err)
	}
	if f.Source() != nil {
		t.Error("Source() should be nil for empty socket path")
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}

func TestFetcher_Source_ReturnsNilForNoop(t *testing.T) {
	f, _ := spire.NewFetcher(context.Background(), spire.Options{
		SocketPath: "unix:///tmp/no-such-spire.sock",
	})
	src := f.Source()
	if src != nil {
		t.Errorf("noop Source() = %v, want nil", src)
	}
}

func TestFetcher_Close_NilSafe(t *testing.T) {
	var f *spire.Fetcher
	err := f.Close()
	if err != nil {
		t.Errorf("nil Fetcher.Close() = %v, want nil", err)
	}
}

func TestFetcher_Source_NilReceiver(t *testing.T) {
	var f *spire.Fetcher
	src := f.Source()
	if src != nil {
		t.Errorf("nil Fetcher.Source() = %v, want nil", src)
	}
}
