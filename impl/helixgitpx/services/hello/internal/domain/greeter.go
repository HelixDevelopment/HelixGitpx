// Package domain holds the business logic for the hello service.
// Pure Go, no framework imports.
package domain

import (
	"context"
	"fmt"

	"github.com/helixgitpx/platform/errors"
	"google.golang.org/grpc/codes"
)

// Counter increments a per-name counter and returns the new value.
type Counter interface {
	Increment(ctx context.Context, name string) (int64, error)
}

// Cache stores the last greeting for a name.
type Cache interface {
	SetLast(ctx context.Context, name, greeting string) error
}

// Emitter publishes a hello.said event.
type Emitter interface {
	Emit(ctx context.Context, name, greeting string, count int64) error
}

// Response is the result of Greet.
type Response struct {
	Greeting string
	Count    int64
}

// Greeter is the hello business-logic aggregate.
type Greeter struct {
	counter Counter
	cache   Cache
	emitter Emitter
}

// NewGreeter constructs a Greeter.
func NewGreeter(c Counter, ca Cache, e Emitter) *Greeter {
	return &Greeter{counter: c, cache: ca, emitter: e}
}

// Greet increments the counter, caches the last greeting, emits an event, and returns the response.
func (g *Greeter) Greet(ctx context.Context, name string) (*Response, error) {
	if name == "" {
		return nil, errors.New(codes.InvalidArgument, "hello", "name is required")
	}
	count, err := g.counter.Increment(ctx, name)
	if err != nil {
		return nil, errors.New(codes.Internal, "hello", "counter increment").Wrap(err)
	}
	greeting := fmt.Sprintf("hello, %s", name)
	if err := g.cache.SetLast(ctx, name, greeting); err != nil {
		return nil, errors.New(codes.Unavailable, "hello", "cache set last").Wrap(err)
	}
	if err := g.emitter.Emit(ctx, name, greeting, count); err != nil {
		return nil, errors.New(codes.Unavailable, "hello", "emit event").Wrap(err)
	}
	return &Response{Greeting: greeting, Count: count}, nil
}
