package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/adapter-pool/internal/adapter"
	"github.com/helixgitpx/helixgitpx/services/adapter-pool/internal/providers/github"
	"github.com/helixgitpx/helixgitpx/services/adapter-pool/internal/providers/gitlab"
)

func setup(t *testing.T) *httptest.Server {
	t.Helper()
	reg := NewRegistry()
	reg.Register(adapter.GitHub, &github.Adapter{})
	reg.Register(adapter.GitLab, &gitlab.Adapter{})
	srv := httptest.NewServer((&Handler{Registry: reg}).Routes())
	t.Cleanup(srv.Close)
	return srv
}

func TestListProviders(t *testing.T) {
	srv := setup(t)
	resp, err := http.Get(srv.URL + "/v1/adapters")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	var body struct {
		Providers []string `json:"providers"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if len(body.Providers) < 2 {
		t.Fatalf("expected at least 2 providers, got %v", body.Providers)
	}
}

func TestProviderHealth_Known(t *testing.T) {
	srv := setup(t)
	resp, err := http.Get(srv.URL + "/v1/adapters/github/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	var body struct {
		Healthy bool `json:"healthy"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if !body.Healthy {
		t.Fatal("github adapter should report healthy")
	}
}

func TestProviderHealth_Unknown(t *testing.T) {
	srv := setup(t)
	resp, err := http.Get(srv.URL + "/v1/adapters/nonexistent/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatalf("want 404 got %d", resp.StatusCode)
	}
}

func TestHealthz(t *testing.T) {
	srv := setup(t)
	resp, err := http.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("want application/json got %q", ct)
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Status != "ok" {
		t.Fatalf("want ok got %q", body.Status)
	}
}
