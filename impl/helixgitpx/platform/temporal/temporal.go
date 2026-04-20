// Package temporal wraps go.temporal.io/sdk with HelixGitpx defaults.
// The stub path (HostPort == "") lets tests and local dev construct a
// Client without requiring a running Temporal frontend.
package temporal

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrUnavailable indicates the Temporal service cannot be reached.
var ErrUnavailable = errors.New("temporal: unavailable")

// Options configures NewClient.
type Options struct {
	HostPort     string
	Namespace    string
	Identity     string
	DialTimeout  time.Duration
	Interceptors []any // placeholder for future opentelemetry/auth interceptors
}

// Dialer abstracts the Temporal SDK's client.Dial so we can swap in a fake
// during tests. A real wiring installs the SDK-backed dialer at package init.
type Dialer func(ctx context.Context, opts Options) (SDKClient, error)

// SDKClient is the narrow surface we need from the Temporal SDK.
// Keeping it narrow lets us mock in unit tests AND lets us vendor the SDK
// only in the binary that actually registers workflows.
type SDKClient interface {
	Close()
}

// Client is a thin wrapper that either delegates to an SDK connection or
// behaves as a no-op (when HostPort is empty).
type Client struct {
	HostPort  string
	Namespace string
	noop      bool
	sdk       SDKClient
}

var activeDialer Dialer

// RegisterDialer installs the real SDK dialer. Services that depend on
// temporal call this in init(); tests can swap it out.
func RegisterDialer(d Dialer) { activeDialer = d }

// NewClient returns a Client. When HostPort is empty, returns a no-op client.
// Otherwise it invokes the registered Dialer, returning ErrUnavailable if
// none is installed.
func NewClient(ctx context.Context, opts Options) (*Client, error) {
	if opts.HostPort == "" {
		return &Client{noop: true, Namespace: opts.Namespace}, nil
	}
	if activeDialer == nil {
		return nil, fmt.Errorf("%w: no dialer registered (call temporal.RegisterDialer)", ErrUnavailable)
	}
	if opts.DialTimeout == 0 {
		opts.DialTimeout = 5 * time.Second
	}
	dctx, cancel := context.WithTimeout(ctx, opts.DialTimeout)
	defer cancel()
	sdk, err := activeDialer(dctx, opts)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	return &Client{
		HostPort:  opts.HostPort,
		Namespace: opts.Namespace,
		sdk:       sdk,
	}, nil
}

// Close releases resources.
func (c *Client) Close() error {
	if c == nil || c.sdk == nil {
		return nil
	}
	c.sdk.Close()
	return nil
}

// IsNoop reports whether the Client was built without a HostPort.
func (c *Client) IsNoop() bool { return c != nil && c.noop }
