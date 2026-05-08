package redis_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	hr "github.com/helixgitpx/platform/redis"
)

func TestKey_AppliesNamespace(t *testing.T) {
	c := hr.Client{Namespace: "hello"}
	got := c.Key("greeting", "world")
	want := "hello:greeting:world"
	if got != want {
		t.Errorf("Key = %q, want %q", got, want)
	}
}

func TestKey_NoNamespace(t *testing.T) {
	c := hr.Client{Namespace: ""}
	got := c.Key("a", "b")
	want := "a:b"
	if got != want {
		t.Errorf("Key = %q, want %q", got, want)
	}
}

func TestIsUnavailable(t *testing.T) {
	if !hr.IsUnavailable(hr.ErrUnavailable) {
		t.Errorf("sentinel not classified")
	}
	if hr.IsUnavailable(errors.New("other")) {
		t.Errorf("other err misclassified")
	}
}

func TestProbe_ReturnsErrorOnNilClient(t *testing.T) {
	probe := hr.Probe(nil)
	err := probe(context.Background())
	if err == nil {
		t.Fatal("expected error for nil client probe")
	}
	if !hr.IsUnavailable(err) {
		t.Errorf("probe(nil) err = %v, want ErrUnavailable", err)
	}
}

func TestOpen_InvalidAddr_ReturnsUnavailable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := hr.Open(ctx, hr.Options{Addr: "localhost:19999"})
	if err == nil {
		t.Fatal("expected error for invalid addr")
	}
	if !hr.IsUnavailable(err) {
		t.Errorf("Open(invalid) err = %v, want ErrUnavailable", err)
	}
}

func TestOpen_RealRedis_Pings(t *testing.T) {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := hr.Open(ctx, hr.Options{
		Addr:      addr,
		Namespace: "antibluff",
	})
	if err != nil {
		t.Skipf("SKIP-OK: #integration — Redis not available at %s: %v", addr, err)
	}
	defer client.Close()

	pingErr := client.Ping(ctx).Err()
	if pingErr != nil {
		t.Fatalf("Ping after Open failed: %v", pingErr)
	}

	key := client.Key("test", "ping")
	got := client.Key("test", "ping")
	if got != "antibluff:test:ping" {
		t.Errorf("Key = %q, want %q", got, "antibluff:test:ping")
	}
	_ = key

	setErr := client.Set(ctx, client.Key("antibluff", "check"), "works", 10*time.Second).Err()
	if setErr != nil {
		t.Fatalf("Set failed: %v", setErr)
	}
	val, getErr := client.Get(ctx, client.Key("antibluff", "check")).Result()
	if getErr != nil {
		t.Fatalf("Get failed: %v", getErr)
	}
	if val != "works" {
		t.Errorf("Get = %q, want %q", val, "works")
	}
}

func TestProbe_RealRedis_ReturnsNil(t *testing.T) {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := hr.Open(ctx, hr.Options{Addr: addr})
	if err != nil {
		t.Skipf("SKIP-OK: #integration — Redis not available at %s: %v", addr, err)
	}
	defer client.Close()

	probe := hr.Probe(client)
	if probeErr := probe(ctx); probeErr != nil {
		t.Errorf("probe(live client) = %v, want nil", probeErr)
	}
}
