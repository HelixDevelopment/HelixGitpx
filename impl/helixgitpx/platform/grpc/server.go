// Package grpc builds a gRPC server preconfigured with the HelixGitpx
// interceptor chain (logging, recovery, telemetry, auth stub).
package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Options configures NewServer.
type Options struct {
	// ServerOptions is appended to the built-in option list.
	ServerOptions []grpc.ServerOption
	// DisableReflection — reflection is registered by default; set true to disable.
	DisableReflection bool
}

// NewServer constructs a *grpc.Server with HelixGitpx defaults:
//   - unary + stream interceptor chain (see interceptor.go)
//   - grpc.health.v1.Health registered and reporting SERVING
//   - server reflection registered (unless disabled)
func NewServer(opts Options) (*grpc.Server, error) {
	so := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryChain()...),
		grpc.ChainStreamInterceptor(streamChain()...),
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
