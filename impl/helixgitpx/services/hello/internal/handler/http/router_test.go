package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/helixgitpx/helixgitpx/services/hello/internal/domain"
	hellohandler "github.com/helixgitpx/helixgitpx/services/hello/internal/handler/http"
)

// unit-test-only fakes (Constitution §II §2 allows mocks in unit tests).
type fakeCounter struct{ n int64 }

func (f *fakeCounter) Increment(_ context.Context, _ string) (int64, error) {
	f.n++
	return f.n, nil
}

type fakeCache struct{}

func (fakeCache) SetLast(_ context.Context, _, _ string) error { return nil }

type fakeEmitter struct{}

func (fakeEmitter) Emit(_ context.Context, _, _ string, _ int64) error { return nil }

func setup() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	hellohandler.Register(r, domain.NewGreeter(&fakeCounter{}, &fakeCache{}, &fakeEmitter{}))
	return r
}

func TestHelloGET_Happy(t *testing.T) {
	srv := httptest.NewServer(setup())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/v1/hello?name=world")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	var body struct {
		Greeting string `json:"greeting"`
		Count    uint64 `json:"count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Greeting == "" {
		t.Fatal("expected non-empty greeting")
	}
	if body.Count == 0 {
		t.Fatal("expected count >= 1 after first request")
	}
}

func TestHelloGET_CounterMonotonic(t *testing.T) {
	srv := httptest.NewServer(setup())
	defer srv.Close()

	fetch := func() uint64 {
		resp, err := http.Get(srv.URL + "/v1/hello?name=alice")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		var body struct {
			Count uint64 `json:"count"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return body.Count
	}
	first := fetch()
	second := fetch()
	if second <= first {
		t.Fatalf("count should be monotonic: first=%d second=%d", first, second)
	}
}
