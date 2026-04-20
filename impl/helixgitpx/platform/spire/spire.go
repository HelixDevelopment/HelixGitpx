// Package spire integrates SPIFFE/SPIRE workload API.
// M1 ships a stub; M2 wires go-spiffe/v2 to fetch SVIDs.
package spire

import (
	"context"
	"errors"
	"os"
)

// ErrUnavailable indicates the SPIRE agent socket is not reachable.
var ErrUnavailable = errors.New("spire: unavailable")

// Options configures NewFetcher.
type Options struct {
	SocketPath string // e.g. "unix:///run/spire/agent.sock"
}

// Fetcher retrieves workload SVIDs. Stubbed in M1.
type Fetcher struct {
	SocketPath string
	noop       bool
}

// NewFetcher returns a Fetcher. When the socket file is absent, returns a no-op fetcher.
//
// TODO(M2): wire github.com/spiffe/go-spiffe/v2/workloadapi.NewX509Source and
// supply SVIDs to grpc/TLS constructors.
func NewFetcher(_ context.Context, opts Options) (*Fetcher, error) {
	if opts.SocketPath == "" {
		return &Fetcher{noop: true}, nil
	}
	if _, err := os.Stat(trimUnix(opts.SocketPath)); err != nil {
		return &Fetcher{noop: true, SocketPath: opts.SocketPath}, nil
	}
	return &Fetcher{SocketPath: opts.SocketPath}, nil
}

// Close releases resources. No-op for M1.
func (f *Fetcher) Close() error { return nil }

func trimUnix(s string) string {
	const prefix = "unix://"
	if len(s) > len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}
