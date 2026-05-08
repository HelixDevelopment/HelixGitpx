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

func TestValidatePush_InvalidJSON(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/pushes/validate", "application/json",
		strings.NewReader(`{bad json`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("want 400 got %d", resp.StatusCode)
	}
	var errBody struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&errBody); err != nil {
		t.Fatalf("decode error body: %v", err)
	}
	if errBody.Code != "invalid_json" {
		t.Fatalf("want code=invalid_json got %q", errBody.Code)
	}
}

func TestValidatePush_ValidBranch_VerifiesContent(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/pushes/validate", "application/json",
		strings.NewReader(`{"ref":"refs/heads/feat/x","repo_id":"r-1","size_bytes":100,"pushes_last_minute":0,"push_limit":10,"max_bytes_per_push":100000}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	var body validateOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if !body.Valid {
		t.Fatalf("expected valid, got error: %q", body.Error)
	}
}

func TestValidatePush_ProtectedMain(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/pushes/validate", "application/json",
		strings.NewReader(`{"ref":"refs/heads/main","repo_id":"r-1","size_bytes":100,"pushes_last_minute":0,"push_limit":10,"max_bytes_per_push":100000}`))
	defer resp.Body.Close()
	var body validateOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if !body.Valid {
		t.Fatalf("main push should be valid: %q", body.Error)
	}
	if !body.Protected {
		t.Fatal("main must be flagged as protected")
	}
}

func TestValidatePush_InvalidRef(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/pushes/validate", "application/json",
		strings.NewReader(`{"ref":"not-a-ref","repo_id":"r-1","size_bytes":100,"pushes_last_minute":0,"push_limit":10,"max_bytes_per_push":100000}`))
	defer resp.Body.Close()
	var body validateOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.Valid {
		t.Fatal("invalid ref should not be valid")
	}
}

func TestValidatePush_QuotaExceeded(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/pushes/validate", "application/json",
		strings.NewReader(`{"ref":"refs/heads/feat/y","repo_id":"r-1","size_bytes":100,"pushes_last_minute":10,"push_limit":10,"max_bytes_per_push":100000}`))
	defer resp.Body.Close()
	var body validateOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.Valid {
		t.Fatal("over quota should not be valid")
	}
}

func TestValidatePush_TooLarge(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/pushes/validate", "application/json",
		strings.NewReader(`{"ref":"refs/heads/feat/y","repo_id":"r-1","size_bytes":1000000,"pushes_last_minute":0,"push_limit":10,"max_bytes_per_push":1000}`))
	defer resp.Body.Close()
	var body validateOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.Valid {
		t.Fatal("oversized push should not be valid")
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
