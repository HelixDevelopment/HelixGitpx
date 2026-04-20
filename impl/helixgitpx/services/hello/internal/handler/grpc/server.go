package grpc

import (
	"context"

	hellopb "github.com/helixgitpx/helixgitpx/gen/go/helixgitpx/hello/v1"
	"github.com/helixgitpx/helixgitpx/services/hello/internal/domain"
)

// Server implements hellopb.HelloServiceServer.
type Server struct {
	hellopb.UnimplementedHelloServiceServer
	Greeter *domain.Greeter
}

// SayHello satisfies hellopb.HelloServiceServer.
func (s *Server) SayHello(ctx context.Context, req *hellopb.SayHelloRequest) (*hellopb.SayHelloResponse, error) {
	r, err := s.Greeter.Greet(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	return &hellopb.SayHelloResponse{Greeting: r.Greeting, Count: r.Count}, nil
}
