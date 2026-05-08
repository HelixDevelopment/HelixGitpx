package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/ai-service/internal/usecase"
)

type fakeLLM struct {
	response string
}

func (f *fakeLLM) Prompt(_ context.Context, _, prompt string) (string, error) {
	return f.response + ":" + prompt, nil
}

func setup(t *testing.T) *httptest.Server {
	t.Helper()
	llm := &fakeLLM{response: "ai-result"}
	uc := &usecase.UseCases{LLM: llm}
	srv := httptest.NewServer(NewHandler(uc).Routes())
	t.Cleanup(srv.Close)
	return srv
}

func TestSummarize(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/ai/summarize", "application/json",
		strings.NewReader(`{"content":"long text here"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	var body promptOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.Result == "" {
		t.Fatal("expected non-empty result")
	}
	if !strings.Contains(body.Result, "Summarize:") {
		t.Fatalf("result should contain prompt prefix: %q", body.Result)
	}
}

func TestConflict(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/ai/conflict", "application/json",
		strings.NewReader(`{"diff":"<<<<<<< HEAD\nfoo\n=======\nbar\n>>>>>>> other"}`))
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	var body promptOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if !strings.Contains(body.Result, "Resolve conflict:") {
		t.Fatalf("result should contain conflict prefix: %q", body.Result)
	}
}

func TestLabels(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/ai/labels", "application/json",
		strings.NewReader(`{"title":"Fix bug","body":"Fixes the crash"}`))
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	var body promptOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if !strings.Contains(body.Result, "Labels for:") {
		t.Fatalf("result should contain label prefix: %q", body.Result)
	}
}

func TestChatOps(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/ai/chatops", "application/json",
		strings.NewReader(`{"prompt":"deploy to prod"}`))
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	var body promptOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if !strings.Contains(body.Result, "deploy to prod") {
		t.Fatalf("result should contain prompt: %q", body.Result)
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
