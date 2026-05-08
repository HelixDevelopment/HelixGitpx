package grpc

import (
	"context"
	"testing"

	hellopb "github.com/helixgitpx/helixgitpx/gen/go/helixgitpx/hello/v1"
	"github.com/helixgitpx/helixgitpx/services/hello/internal/domain"
)

type fakeCounter struct{ n int64 }

func (f *fakeCounter) Increment(_ context.Context, _ string) (int64, error) {
	f.n++
	return f.n, nil
}

type fakeCache struct{}

func (f *fakeCache) SetLast(_ context.Context, _, _ string) error { return nil }

type fakeEmitter struct{ emitted bool }

func (f *fakeEmitter) Emit(_ context.Context, _, _ string, _ int64) error {
	f.emitted = true
	return nil
}

func TestServer_SayHello(t *testing.T) {
	counter := &fakeCounter{}
	emitter := &fakeEmitter{}
	greeter := domain.NewGreeter(counter, &fakeCache{}, emitter)

	srv := &Server{Greeter: greeter}
	resp, err := srv.SayHello(context.Background(), &hellopb.SayHelloRequest{Name: "world"})
	if err != nil {
		t.Fatalf("SayHello: %v", err)
	}
	if resp.Greeting != "hello, world" {
		t.Errorf("Greeting = %q, want %q", resp.Greeting, "hello, world")
	}
	if resp.Count != 1 {
		t.Errorf("Count = %d, want 1", resp.Count)
	}
	if !emitter.emitted {
		t.Error("expected Emit to be called")
	}
}

func TestServer_SayHello_IncrementsCounter(t *testing.T) {
	counter := &fakeCounter{}
	greeter := domain.NewGreeter(counter, &fakeCache{}, &fakeEmitter{})
	srv := &Server{Greeter: greeter}

	resp1, _ := srv.SayHello(context.Background(), &hellopb.SayHelloRequest{Name: "a"})
	resp2, _ := srv.SayHello(context.Background(), &hellopb.SayHelloRequest{Name: "b"})

	if resp1.Count >= resp2.Count {
		t.Errorf("Count not monotonically increasing: %d then %d", resp1.Count, resp2.Count)
	}
}

func TestServer_SayHello_EmptyName(t *testing.T) {
	greeter := domain.NewGreeter(&fakeCounter{}, &fakeCache{}, &fakeEmitter{})
	srv := &Server{Greeter: greeter}

	_, err := srv.SayHello(context.Background(), &hellopb.SayHelloRequest{Name: ""})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}
