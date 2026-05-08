package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/helixgitpx/platform/log"
)

func freePort(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return fmt.Sprintf("127.0.0.1:%d", port)
}

func TestRun_ServesHealthzAndShutsDown(t *testing.T) {
	addr := freePort(t)
	t.Setenv("LIVE_EVENTS_HTTP_ADDR", addr)

	ctx, cancel := context.WithCancel(context.Background())
	lg := log.Default()

	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, lg)
	}()

	var resp *http.Response
	var err error
	for i := 0; i < 50; i++ {
		resp, err = http.Get("http://" + addr + "/healthz")
		if err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("server never became ready: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("healthz: want 200 got %d", resp.StatusCode)
	}

	resp2, err := http.Get("http://" + addr + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&body); err != nil {
		t.Fatalf("decode healthz: %v", err)
	}
	if body.Status != "ok" {
		t.Fatalf("want status=ok got %q", body.Status)
	}

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not shut down within 5s")
	}
}

func TestEnvOrDefault_Empty(t *testing.T) {
	if got := envOrDefault("NOT_SET_XYZ", ":8080"); got != ":8080" {
		t.Fatalf("want :8080 got %q", got)
	}
}

func TestEnvOrDefault_Set(t *testing.T) {
	t.Setenv("TEST_LIVE_EVENTS_KEY", "  :9999  ")
	if got := envOrDefault("TEST_LIVE_EVENTS_KEY", ":8080"); got != ":9999" {
		t.Fatalf("want :9999 got %q", got)
	}
}
