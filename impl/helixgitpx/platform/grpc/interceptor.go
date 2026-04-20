package grpc

import (
	"context"
	"runtime/debug"

	"github.com/helixgitpx/platform/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func unaryChain() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		recoveryUnary,
		loggingUnary,
	}
}

func streamChain() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		recoveryStream,
	}
}

func recoveryUnary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.FromContext(ctx).Error("panic in grpc handler",
				"method", info.FullMethod, "panic", r, "stack", string(debug.Stack()))
			err = status.Errorf(codes.Internal, "internal error")
		}
	}()
	return handler(ctx, req)
}

func recoveryStream(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.FromContext(ss.Context()).Error("panic in grpc stream",
				"method", info.FullMethod, "panic", r, "stack", string(debug.Stack()))
			err = status.Errorf(codes.Internal, "internal error")
		}
	}()
	return handler(srv, ss)
}

func loggingUnary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)
	lg := log.FromContext(ctx).With("method", info.FullMethod)
	if err != nil {
		lg.Error("grpc call failed", "err", err.Error())
	} else {
		lg.Debug("grpc call ok")
	}
	return resp, err
}
