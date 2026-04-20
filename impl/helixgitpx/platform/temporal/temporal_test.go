package temporal_test

import (
	"context"
	"errors"
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
	if !c.IsNoop() {
		t.Fatalf("expected no-op client")
	}
	if err := c.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}

type fakeSDK struct{ closed bool }

func (f *fakeSDK) Close() { f.closed = true }

func TestNewClient_ErrorsWhenNoDialerRegistered(t *testing.T) {
	temporal.RegisterDialer(nil)
	_, err := temporal.NewClient(context.Background(), temporal.Options{HostPort: "localhost:7233"})
	if !errors.Is(err, temporal.ErrUnavailable) {
		t.Fatalf("want ErrUnavailable, got %v", err)
	}
}

func TestNewClient_UsesRegisteredDialer(t *testing.T) {
	fake := &fakeSDK{}
	temporal.RegisterDialer(func(_ context.Context, _ temporal.Options) (temporal.SDKClient, error) {
		return fake, nil
	})
	t.Cleanup(func() { temporal.RegisterDialer(nil) })

	c, err := temporal.NewClient(context.Background(), temporal.Options{HostPort: "localhost:7233", Namespace: "default"})
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	if c.IsNoop() {
		t.Fatal("should not be no-op")
	}
	if err := c.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if !fake.closed {
		t.Fatal("SDK not closed")
	}
}
