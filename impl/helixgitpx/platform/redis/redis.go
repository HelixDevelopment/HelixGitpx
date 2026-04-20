// Package redis wraps go-redis v9 with namespaced keys and typed unavailability.
package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

// ErrUnavailable is returned when the server cannot be reached.
var ErrUnavailable = errors.New("redis: unavailable")

// Options configures Open.
type Options struct {
	Addr      string
	Password  string
	DB        int
	Namespace string
}

// Client wraps *redis.Client with namespace helpers.
type Client struct {
	*redis.Client
	Namespace string
}

// Open constructs a Client and pings the server.
func Open(ctx context.Context, opts Options) (*Client, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	})
	if err := rc.Ping(ctx).Err(); err != nil {
		_ = rc.Close()
		return nil, errors.Join(ErrUnavailable, err)
	}
	return &Client{Client: rc, Namespace: opts.Namespace}, nil
}

// Key joins parts with ":" and prepends the namespace.
func (c Client) Key(parts ...string) string {
	if c.Namespace == "" {
		return strings.Join(parts, ":")
	}
	return fmt.Sprintf("%s:%s", c.Namespace, strings.Join(parts, ":"))
}

// IsUnavailable reports whether err wraps ErrUnavailable.
func IsUnavailable(err error) bool { return errors.Is(err, ErrUnavailable) }

// Probe returns a health probe function.
func Probe(c *Client) func(context.Context) error {
	return func(ctx context.Context) error {
		if c == nil || c.Client == nil {
			return ErrUnavailable
		}
		return c.Ping(ctx).Err()
	}
}
