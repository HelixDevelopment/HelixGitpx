package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func setup(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer((&Handler{}).Routes())
	t.Cleanup(srv.Close)
	return srv
}

func testInvalidJSON(t *testing.T, srv *httptest.Server, path string) {
	t.Helper()
	resp, err := http.Post(srv.URL+path, "application/json",
		strings.NewReader(`{bad`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("%s: want 400 got %d", path, resp.StatusCode)
	}
	var errBody struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&errBody); err != nil {
		t.Fatalf("decode error body: %v", err)
	}
	if errBody.Code != "invalid_json" {
		t.Fatalf("%s: want code=invalid_json got %q", path, errBody.Code)
	}
}

func TestEncode_InvalidJSON(t *testing.T) {
	srv := setup(t)
	testInvalidJSON(t, srv, "/v1/tokens/encode")
}

func TestDecode_InvalidJSON(t *testing.T) {
	srv := setup(t)
	testInvalidJSON(t, srv, "/v1/tokens/decode")
}

func TestMatch_InvalidJSON(t *testing.T) {
	srv := setup(t)
	testInvalidJSON(t, srv, "/v1/events/match")
}

func TestEncodeDecodeToken_RoundTrip(t *testing.T) {
	srv := setup(t)
	now := time.Now().Unix()
	resp, err := http.Post(srv.URL+"/v1/tokens/encode", "application/json",
		strings.NewReader(fmt.Sprintf(`{"offset":42,"timestamp":%d}`, now)))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("encode: want 200 got %d", resp.StatusCode)
	}
	var enc encodeOut
	_ = json.NewDecoder(resp.Body).Decode(&enc)
	if enc.Token == "" {
		t.Fatal("expected non-empty token")
	}

	resp2, err := http.Post(srv.URL+"/v1/tokens/decode", "application/json",
		strings.NewReader(fmt.Sprintf(`{"token":"%s","retention_seconds":999999999}`, enc.Token)))
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != 200 {
		t.Fatalf("decode: want 200 got %d", resp2.StatusCode)
	}
	var dec decodeOut
	_ = json.NewDecoder(resp2.Body).Decode(&dec)
	if dec.Offset != 42 {
		t.Fatalf("want offset 42 got %d", dec.Offset)
	}
	if dec.Error != "" {
		t.Fatalf("unexpected error: %q", dec.Error)
	}
}

func TestDecodeToken_Empty(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/tokens/decode", "application/json",
		strings.NewReader(`{"token":"","retention_seconds":3600}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var dec decodeOut
	_ = json.NewDecoder(resp.Body).Decode(&dec)
	if dec.Error == "" {
		t.Fatal("expected error for empty token")
	}
}

func TestDecodeToken_Malformed(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/tokens/decode", "application/json",
		strings.NewReader(`{"token":"not-base64!!!","retention_seconds":3600}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var dec decodeOut
	_ = json.NewDecoder(resp.Body).Decode(&dec)
	if dec.Error == "" {
		t.Fatal("expected error for malformed token")
	}
}

func TestMatchEvent_AllowAll(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/events/match", "application/json",
		strings.NewReader(`{"sub_repos":null,"sub_types":null,"repo_id":"r1","event_type":"PUSH"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var body matchOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if !body.Match {
		t.Fatal("nil filters should match everything")
	}
}

func TestMatchEvent_FilterByRepo(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/events/match", "application/json",
		strings.NewReader(`{"sub_repos":["r1","r2"],"sub_types":null,"repo_id":"r1","event_type":"PUSH"}`))
	defer resp.Body.Close()
	var body matchOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if !body.Match {
		t.Fatal("repo in filter should match")
	}
}

func TestMatchEvent_FilterReject(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/events/match", "application/json",
		strings.NewReader(`{"sub_repos":["r1"],"sub_types":null,"repo_id":"r3","event_type":"PUSH"}`))
	defer resp.Body.Close()
	var body matchOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.Match {
		t.Fatal("repo not in filter should not match")
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
