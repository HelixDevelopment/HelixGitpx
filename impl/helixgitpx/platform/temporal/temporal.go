// Package temporal wraps go.temporal.io/sdk with HelixGitpx defaults.
// M1 ships a no-op stub; M5 wires real workflow/activity registration.
package temporal

import (
	"context"
	"errors"
)

// ErrUnavailable indicates the Temporal service cannot be reached.
var ErrUnavailable = errors.New("temporal: unavailable")

// Options configures NewClient.
type Options struct {
	HostPort  string
	Namespace string
}

// Client is a placeholder for the real Temporal client wired in M5.
type Client struct {
	HostPort  string
	Namespace string
	noop      bool
}

// NewClient returns a Client. When HostPort is empty, returns a no-op client.
//
// TODO(M5): wire go.temporal.io/sdk/client.Dial; register workers and workflows.
func NewClient(_ context.Context, opts Options) (*Client, error) {
	if opts.HostPort == "" {
		return &Client{noop: true}, nil
	}
	return &Client{HostPort: opts.HostPort, Namespace: opts.Namespace}, nil
}

// Close releases resources. No-op for M1.
func (c *Client) Close() error { return nil }
