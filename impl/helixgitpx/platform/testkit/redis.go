package testkit

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/redis"
)

// StartRedis launches a Redis 7 container and returns "host:port".
func StartRedis(t testing.TB) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	ctr, err := redis.Run(ctx, "redis:7-alpine")
	if err != nil {
		t.Fatalf("testkit.StartRedis: %v", err)
	}
	ep, err := ctr.Endpoint(ctx, "")
	if err != nil {
		t.Fatalf("testkit.StartRedis endpoint: %v", err)
	}
	t.Cleanup(func() { _ = ctr.Terminate(context.Background()) })
	_ = time.Now
	return ep
}
