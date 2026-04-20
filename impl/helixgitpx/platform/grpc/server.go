// Package grpc builds a gRPC server preconfigured with the HelixGitpx
// interceptor chain (logging, recovery, telemetry, auth stub).
package grpc

import (
	"github.com/helixgitpx/platform/spire"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Options configures NewServer.
type Options struct {
	ServerOptions     []grpc.ServerOption
	DisableReflection bool
	// Fetcher — when non-nil and Fetcher.Source() returns non-nil, the
	// server terminates mTLS with the workload's X.509 SVID and accepts
	// any client SVID in the same trust domain.
	Fetcher *spire.Fetcher
}

// NewServer constructs a *grpc.Server with HelixGitpx defaults:
//   - unary + stream interceptor chain (see interceptor.go)
//   - grpc.health.v1.Health registered and reporting SERVING
//   - server reflection registered (unless disabled)
//   - optional SVID-backed mTLS when Fetcher is present
func NewServer(opts Options) (*grpc.Server, error) {
	so := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryChain()...),
		grpc.ChainStreamInterceptor(streamChain()...),
	}

	if opts.Fetcher != nil {
		if src := opts.Fetcher.Source(); src != nil {
			tlsCfg := tlsconfig.MTLSServerConfig(src, src, tlsconfig.AuthorizeAny())
			so = append(so, grpc.Creds(credentials.NewTLS(tlsCfg)))
		}
	}

	so = append(so, opts.ServerOptions...)

	s := grpc.NewServer(so...)

	hs := health.NewServer()
	hs.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, hs)

	if !opts.DisableReflection {
		reflection.Register(s)
	}
	return s, nil
}
