package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setup(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer((&Handler{}).Routes())
	t.Cleanup(srv.Close)
	return srv
}

func TestClassify_Transient(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/retry/classify", "application/json",
		strings.NewReader(`{"http_status":503}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	var body classifyOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.Kind != "transient" {
		t.Fatalf("want transient got %q", body.Kind)
	}
}

func TestClassify_RateLimit(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/retry/classify", "application/json",
		strings.NewReader(`{"http_status":429}`))
	defer resp.Body.Close()
	var body classifyOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.Kind != "rate_limit" {
		t.Fatalf("want rate_limit got %q", body.Kind)
	}
}

func TestClassify_AuthFailed(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/retry/classify", "application/json",
		strings.NewReader(`{"http_status":401}`))
	defer resp.Body.Close()
	var body classifyOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.Kind != "auth_failed" {
		t.Fatalf("want auth_failed got %q", body.Kind)
	}
}

func TestClassify_Permanent(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/retry/classify", "application/json",
		strings.NewReader(`{"http_status":0,"permanent":true}`))
	defer resp.Body.Close()
	var body classifyOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.Kind != "permanent" {
		t.Fatalf("want permanent got %q", body.Kind)
	}
}

func TestBackoff(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/retry/backoff", "application/json",
		strings.NewReader(`{"attempt":1,"base_ms":1000,"max_ms":60000}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var body backoffOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.DelayMs != 1000 {
		t.Fatalf("attempt 1: want 1000ms got %d", body.DelayMs)
	}
}

func TestBackoff_Exponential(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/retry/backoff", "application/json",
		strings.NewReader(`{"attempt":4,"base_ms":1000,"max_ms":60000}`))
	defer resp.Body.Close()
	var body backoffOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.DelayMs != 8000 {
		t.Fatalf("attempt 4: want 8000ms got %d", body.DelayMs)
	}
}

func TestBackoff_Cap(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/retry/backoff", "application/json",
		strings.NewReader(`{"attempt":10,"base_ms":1000,"max_ms":5000}`))
	defer resp.Body.Close()
	var body backoffOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.DelayMs > 5000 {
		t.Fatalf("should be capped at 5000ms, got %d", body.DelayMs)
	}
}

func TestDecision_TransientRetries(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/retry/decision", "application/json",
		strings.NewReader(`{"http_status":503,"attempt":1,"max_attempts":5}`))
	defer resp.Body.Close()
	var body decisionOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if !body.ShouldRetry {
		t.Fatal("transient should retry")
	}
	if body.GoesToDLQ {
		t.Fatal("transient below cap should not go to DLQ")
	}
}

func TestDecision_PermanentDLQ(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/retry/decision", "application/json",
		strings.NewReader(`{"http_status":0,"permanent":true,"attempt":1,"max_attempts":5}`))
	defer resp.Body.Close()
	var body decisionOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.ShouldRetry {
		t.Fatal("permanent should not retry")
	}
	if !body.GoesToDLQ {
		t.Fatal("permanent should go to DLQ")
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
