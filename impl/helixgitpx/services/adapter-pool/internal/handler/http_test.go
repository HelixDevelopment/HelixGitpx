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
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("want application/json got %q", ct)
	}
	var body struct {
		Providers []string `json:"providers"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body.Providers) < 2 {
		t.Fatalf("expected at least 2 providers, got %v", body.Providers)
	}
	have := map[string]bool{}
	for _, p := range body.Providers {
		have[p] = true
	}
	if !have["github"] {
		t.Fatal("providers list should contain 'github'")
	}
	if !have["gitlab"] {
		t.Fatal("providers list should contain 'gitlab'")
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
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("want application/json got %q", ct)
	}
	var body struct {
		Provider      string `json:"provider"`
		Healthy       bool   `json:"healthy"`
		DefaultBranch string `json:"default_branch"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Provider != "github" {
		t.Fatalf("want provider=github got %q", body.Provider)
	}
	if !body.Healthy {
		t.Fatal("github adapter should report healthy")
	}
	if body.DefaultBranch != "main" {
		t.Fatalf("want default_branch=main got %q", body.DefaultBranch)
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
	var errBody struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&errBody); err != nil {
		t.Fatalf("decode error body: %v", err)
	}
	if errBody.Code != "not_found" {
		t.Fatalf("want code=not_found got %q", errBody.Code)
	}
}

func TestProviderHealth_GitLab(t *testing.T) {
	srv := setup(t)
	resp, err := http.Get(srv.URL + "/v1/adapters/gitlab/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	var body struct {
		Provider string `json:"provider"`
		Healthy  bool   `json:"healthy"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Provider != "gitlab" {
		t.Fatalf("want provider=gitlab got %q", body.Provider)
	}
	if !body.Healthy {
		t.Fatal("gitlab adapter should report healthy")
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
