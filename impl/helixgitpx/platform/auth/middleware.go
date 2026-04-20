package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type userIDKey struct{}

// UserIDFromContext retrieves the user id set by the interceptor, or "" if absent.
func UserIDFromContext(ctx context.Context) (string, bool) {
	s, ok := ctx.Value(userIDKey{}).(string)
	return s, ok
}

// UnaryInterceptor returns a gRPC unary interceptor that validates the
// bearer token from metadata, injects the user id into context, and
// rejects unauthenticated calls with codes.Unauthenticated.
func UnaryInterceptor(v *Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "no metadata")
		}
		authz := md.Get("authorization")
		if len(authz) == 0 {
			return nil, status.Error(codes.Unauthenticated, "no authorization header")
		}
		token := strings.TrimPrefix(authz[0], "Bearer ")
		claims, err := v.Validate(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}
		ctx = context.WithValue(ctx, userIDKey{}, claims.Subject)
		return handler(ctx, req)
	}
}
