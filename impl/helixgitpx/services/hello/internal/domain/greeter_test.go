package domain_test

import (
	"context"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/hello/internal/domain"
)

type fakeCounter struct {
	count int64
	err   error
}

func (f *fakeCounter) Increment(_ context.Context, name string) (int64, error) {
	if f.err != nil {
		return 0, f.err
	}
	f.count++
	return f.count, nil
}

type fakeCache struct {
	last string
}

func (f *fakeCache) SetLast(_ context.Context, name, greeting string) error {
	f.last = greeting
	return nil
}

type fakeEmitter struct {
	events int
}

func (f *fakeEmitter) Emit(_ context.Context, name, greeting string, count int64) error {
	f.events++
	return nil
}

func TestGreeter_Greet_ReturnsFormattedGreeting(t *testing.T) {
	g := domain.NewGreeter(&fakeCounter{}, &fakeCache{}, &fakeEmitter{})
	resp, err := g.Greet(context.Background(), "world")
	if err != nil {
		t.Fatalf("Greet: %v", err)
	}
	if resp.Greeting != "hello, world" {
		t.Errorf("Greeting = %q", resp.Greeting)
	}
	if resp.Count != 1 {
		t.Errorf("Count = %d, want 1", resp.Count)
	}
}

func TestGreeter_Greet_EmptyNameFails(t *testing.T) {
	g := domain.NewGreeter(&fakeCounter{}, &fakeCache{}, &fakeEmitter{})
	_, err := g.Greet(context.Background(), "")
	if err == nil {
		t.Fatalf("expected error for empty name")
	}
}
