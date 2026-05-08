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
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("want application/json content-type, got %q", ct)
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
	if body.ElapsedMs < 0 {
		t.Fatalf("want elapsed_ms >= 0 got %d", body.ElapsedMs)
	}
	if len(body.Engines) != 2 || body.Engines[0] != "meili" || body.Engines[1] != "qdrant" {
		t.Fatalf("want engines [meili qdrant] got %v", body.Engines)
	}
	if body.Hits[0].Score <= 0 {
		t.Fatalf("top hit should have positive score, got %f", body.Hits[0].Score)
	}
	if len(body.Hits[0].PerEngine) == 0 {
		t.Fatal("top hit should have per_engine scores")
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

func TestSearch_EmptyQueryReturnsEmptyHits(t *testing.T) {
	h := &Handler{
		Engines: []engines.Engine{
			&fakeEngine{name: "meili", hits: []engines.Hit{}},
		},
	}
	srv := httptest.NewServer(h.Routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/v1/search")
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
	if len(body.Hits) != 0 {
		t.Fatalf("want 0 hits for empty query got %d", len(body.Hits))
	}
}

func TestSearch_LimitParamTruncatesResults(t *testing.T) {
	h := &Handler{
		Engines: []engines.Engine{
			&fakeEngine{name: "meili", hits: []engines.Hit{{ID: "a"}, {ID: "b"}, {ID: "c"}}},
		},
	}
	srv := httptest.NewServer(h.Routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/v1/search?q=test&limit=2")
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
	if len(body.Hits) != 2 {
		t.Fatalf("limit=2 should truncate to 2 hits, got %d", len(body.Hits))
	}
}

func TestSearch_NoEnginesReturnsEmpty(t *testing.T) {
	h := &Handler{Engines: nil}
	srv := httptest.NewServer(h.Routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/v1/search?q=test")
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
	if len(body.Hits) != 0 {
		t.Fatalf("want 0 hits with no engines got %d", len(body.Hits))
	}
	if len(body.Engines) != 0 {
		t.Fatalf("want 0 engines in response got %v", body.Engines)
	}
}

func TestHealthz_ReturnsStatusOK(t *testing.T) {
	h := &Handler{Engines: nil}
	srv := httptest.NewServer(h.Routes())
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("want application/json content-type, got %q", ct)
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Status != "ok" {
		t.Fatalf("want status=ok got %q", body.Status)
	}
}
