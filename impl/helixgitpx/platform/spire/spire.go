// Package spire integrates SPIFFE/SPIRE workload API via go-spiffe/v2.
// M2 lifts the M1 no-op stub; a non-empty SocketPath now returns a live
// fetcher that streams X.509 SVIDs from the SPIRE agent.
package spire

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

// ErrUnavailable indicates the SPIRE agent socket is not reachable.
var ErrUnavailable = errors.New("spire: unavailable")

// Options configures NewFetcher.
type Options struct {
	SocketPath string // e.g. "unix:///run/spire/agent.sock"
}

// Fetcher streams workload SVIDs from the SPIRE agent.
type Fetcher struct {
	source *workloadapi.X509Source
	noop   bool
}

// NewFetcher returns a Fetcher. When SocketPath is empty or the socket
// file is absent, returns a no-op fetcher so callers can remain agnostic
// on dev machines without SPIRE.
func NewFetcher(ctx context.Context, opts Options) (*Fetcher, error) {
	if opts.SocketPath == "" {
		return &Fetcher{noop: true}, nil
	}
	if _, err := os.Stat(trimUnix(opts.SocketPath)); err != nil {
		return &Fetcher{noop: true}, nil
	}
	src, err := workloadapi.NewX509Source(ctx,
		workloadapi.WithClientOptions(workloadapi.WithAddr(opts.SocketPath)),
	)
	if err != nil {
		return nil, fmt.Errorf("spire: X509Source: %w", errors.Join(ErrUnavailable, err))
	}
	return &Fetcher{source: src}, nil
}

// Source returns the underlying *workloadapi.X509Source, or nil when no-op.
// Callers pass this to grpc.Creds / tls.Config builders.
func (f *Fetcher) Source() *workloadapi.X509Source {
	if f == nil || f.noop {
		return nil
	}
	return f.source
}

// Close releases resources.
func (f *Fetcher) Close() error {
	if f == nil || f.source == nil {
		return nil
	}
	return f.source.Close()
}

func trimUnix(s string) string {
	const prefix = "unix://"
	if len(s) > len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}
