package repo

import (
	"context"
	"os"
	"testing"
	"time"

	hr "github.com/helixgitpx/platform/redis"
)

func TestCacheRedis_SetLast_RealRedis(t *testing.T) {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	ctx := context.Background()

	client, err := hr.Open(ctx, hr.Options{Addr: addr, Namespace: "test_hello"})
	if err != nil {
		t.Skipf("SKIP-OK: #integration — Redis not available: %v", err)
	}
	defer client.Close()

	cache := &CacheRedis{Client: client, TTL: 30 * time.Second}

	if err := cache.SetLast(ctx, "alice", "hello, alice"); err != nil {
		t.Fatalf("SetLast: %v", err)
	}

	key := client.Key("last", "alice")
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		t.Fatalf("Get after SetLast: %v", err)
	}
	if val != "hello, alice" {
		t.Errorf("got %q, want %q", val, "hello, alice")
	}
}

func TestCacheRedis_SetLast_DefaultTTL(t *testing.T) {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	ctx := context.Background()

	client, err := hr.Open(ctx, hr.Options{Addr: addr, Namespace: "test_hello"})
	if err != nil {
		t.Skipf("SKIP-OK: #integration — Redis not available: %v", err)
	}
	defer client.Close()

	cache := &CacheRedis{Client: client} // TTL=0, should default to 10m

	if err := cache.SetLast(ctx, "bob", "hello, bob"); err != nil {
		t.Fatalf("SetLast with default TTL: %v", err)
	}

	val, err := client.Get(ctx, client.Key("last", "bob")).Result()
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if val != "hello, bob" {
		t.Errorf("got %q, want %q", val, "hello, bob")
	}
}
