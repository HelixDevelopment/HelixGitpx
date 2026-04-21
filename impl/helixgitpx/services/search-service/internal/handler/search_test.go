package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/helixgitpx/helixgitpx/services/search-service/internal/engines"
)

// fakeEngine is allowed here: we're in a UNIT test (Constitution §II §2).
type fakeEngine struct {
	name  string
	hits  []engines.Hit
	err   error
	delay time.Duration
}

func (f *fakeEngine) Name() string { return f.name }
func (f *fakeEngine) Search(ctx context.Context, _ engines.Query) ([]engines.Hit, error) {
	if f.delay > 0 {
		select {
		case <-time.After(f.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return f.hits, f.err
}

func TestSearch_FusesHitsFromMultipleEngines(t *testing.T) {
	h := &Handler{
		Engines: []engines.Engine{
			&fakeEngine{name: "meili", hits: []engines.Hit{{ID: "a"}, {ID: "b"}, {ID: "c"}}},
			&fakeEngine{name: "qdrant", hits: []engines.Hit{{ID: "b"}, {ID: "d"}}},
		},
	}
	srv := httptest.NewServer(h.Routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/v1/search?q=hello&limit=10")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	var body response
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body.Hits) != 4 {
		t.Fatalf("want 4 unique hits, got %d", len(body.Hits))
	}
	if body.Hits[0].ID != "b" {
		t.Fatalf("b appears in both engines, should rank first; got %s", body.Hits[0].ID)
	}
}

func TestSearch_TolerantOfFailingEngine(t *testing.T) {
	h := &Handler{
		Engines: []engines.Engine{
			&fakeEngine{name: "good", hits: []engines.Hit{{ID: "x"}}},
			&fakeEngine{name: "broken", err: context.DeadlineExceeded},
		},
	}
	srv := httptest.NewServer(h.Routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/v1/search?q=hi")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("partial failure must still return 200, got %d", resp.StatusCode)
	}
	var body response
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if len(body.Hits) != 1 || body.Hits[0].ID != "x" {
		t.Fatalf("expected the good engine's result to still show; got %+v", body.Hits)
	}
}

func TestHealthz(t *testing.T) {
	h := &Handler{Engines: nil}
	srv := httptest.NewServer(h.Routes())
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatal("healthz must be 200")
	}
}
