package handler

import (
	"context"
	"encoding/json"
	"fmt"
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

type errLLM struct{}

func (e *errLLM) Prompt(_ context.Context, _, _ string) (string, error) {
	return "", fmt.Errorf("llm unavailable")
}

func setup(t *testing.T) *httptest.Server {
	t.Helper()
	llm := &fakeLLM{response: "ai-result"}
	uc := &usecase.UseCases{LLM: llm}
	srv := httptest.NewServer(NewHandler(uc).Routes())
	t.Cleanup(srv.Close)
	return srv
}

func setupWithLLM(llm usecase.Client) *httptest.Server {
	uc := &usecase.UseCases{LLM: llm}
	srv := httptest.NewServer(NewHandler(uc).Routes())
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

func TestSummarize_InvalidJSON(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/ai/summarize", "application/json",
		strings.NewReader(`{bad json`))
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("want 400 got %d", resp.StatusCode)
	}
	var errBody struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&errBody)
	if errBody.Code != "invalid_json" {
		t.Fatalf("want code=invalid_json got %q", errBody.Code)
	}
}

func TestConflict_InvalidJSON(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/ai/conflict", "application/json",
		strings.NewReader(`not json at all`))
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("want 400 got %d", resp.StatusCode)
	}
}

func TestLabels_InvalidJSON(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/ai/labels", "application/json",
		strings.NewReader(``))
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("want 400 got %d", resp.StatusCode)
	}
}

func TestChatOps_InvalidJSON(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/ai/chatops", "application/json",
		strings.NewReader(`{`))
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("want 400 got %d", resp.StatusCode)
	}
}

func TestSummarize_LLMError(t *testing.T) {
	srv := setupWithLLM(&errLLM{})
	defer srv.Close()
	resp, _ := http.Post(srv.URL+"/v1/ai/summarize", "application/json",
		strings.NewReader(`{"content":"test"}`))
	defer resp.Body.Close()
	if resp.StatusCode != 500 {
		t.Fatalf("want 500 got %d", resp.StatusCode)
	}
	var errBody struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&errBody)
	if errBody.Code != "summarize_failed" {
		t.Fatalf("want code=summarize_failed got %q", errBody.Code)
	}
	if errBody.Message != "llm unavailable" {
		t.Fatalf("want message='llm unavailable' got %q", errBody.Message)
	}
}

func TestConflict_LLMError(t *testing.T) {
	srv := setupWithLLM(&errLLM{})
	defer srv.Close()
	resp, _ := http.Post(srv.URL+"/v1/ai/conflict", "application/json",
		strings.NewReader(`{"diff":"test diff"}`))
	defer resp.Body.Close()
	if resp.StatusCode != 500 {
		t.Fatalf("want 500 got %d", resp.StatusCode)
	}
}

func TestSummarize_PromptConstruction(t *testing.T) {
	srv := setup(t)
	content := "This is a long piece of text about HelixGitpx."
	resp, _ := http.Post(srv.URL+"/v1/ai/summarize", "application/json",
		strings.NewReader(`{"content":"`+content+`"}`))
	defer resp.Body.Close()
	var body promptOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if !strings.Contains(body.Result, content) {
		t.Fatalf("result should contain original content: %q", body.Result)
	}
	if !strings.HasPrefix(body.Result, "ai-result:Summarize: ") {
		t.Fatalf("result should start with 'ai-result:Summarize: ': %q", body.Result)
	}
}

func TestConflict_PromptConstruction(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/ai/conflict", "application/json",
		strings.NewReader(`{"diff":"merge conflict here"}`))
	defer resp.Body.Close()
	var body promptOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if !strings.Contains(body.Result, "Resolve conflict: merge conflict here") {
		t.Fatalf("result should contain conflict prompt: %q", body.Result)
	}
}

func TestLabels_PromptConstruction(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/ai/labels", "application/json",
		strings.NewReader(`{"title":"Bug fix","body":"Fixes crash"}`))
	defer resp.Body.Close()
	var body promptOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if !strings.Contains(body.Result, "Labels for: Bug fix") {
		t.Fatalf("result should contain label prefix with title: %q", body.Result)
	}
	if !strings.Contains(body.Result, "Fixes crash") {
		t.Fatalf("result should contain body: %q", body.Result)
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
