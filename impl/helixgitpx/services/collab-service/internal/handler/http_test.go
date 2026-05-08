package handler

import (
	"encoding/json"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setup(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(NewHandler().Routes())
	t.Cleanup(srv.Close)
	return srv
}

func TestValidateDoc_Valid(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/docs/validate", "application/json",
		strings.NewReader(`{"doc_id":"d1","actor_id":"alice"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	var body validateDocOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if !body.Valid {
		t.Fatalf("expected valid, got error: %q", body.Error)
	}
}

func TestValidateDoc_EmptyDoc(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/docs/validate", "application/json",
		strings.NewReader(`{"doc_id":"  ","actor_id":"alice"}`))
	defer resp.Body.Close()
	var body validateDocOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.Valid {
		t.Fatal("empty doc should not be valid")
	}
}

func TestValidateDoc_EmptyActor(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/docs/validate", "application/json",
		strings.NewReader(`{"doc_id":"d1","actor_id":""}`))
	defer resp.Body.Close()
	var body validateDocOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.Valid {
		t.Fatal("empty actor should not be valid")
	}
}

func TestSnapshotCheck_Allowed(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/docs/snapshot-check", "application/json",
		strings.NewReader(`{"snapshot":"dG9vLXNtYWxs"}`))
	defer resp.Body.Close()
	var body snapshotCheckOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if !body.Allowed {
		t.Fatalf("small snapshot should be allowed: %q", body.Error)
	}
}

func TestSnapshotCheck_TooLarge(t *testing.T) {
	srv := setup(t)
	raw := make([]byte, 9*1024*1024)
	for i := range raw {
		raw[i] = byte(i % 256)
	}
	encoded := base64.StdEncoding.EncodeToString(raw)
	payload := `{"snapshot":"` + encoded + `"}`
	resp, _ := http.Post(srv.URL+"/v1/docs/snapshot-check", "application/json",
		strings.NewReader(payload))
	defer resp.Body.Close()
	var body snapshotCheckOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.Allowed {
		t.Fatal("oversized snapshot should not be allowed")
	}
}

func TestParticipantCheck_Allowed(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/docs/participant-check", "application/json",
		strings.NewReader(`{"current_participants":10}`))
	defer resp.Body.Close()
	var body participantCheckOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if !body.Allowed {
		t.Fatalf("below cap should be allowed: %q", body.Error)
	}
}

func TestParticipantCheck_AtCap(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/docs/participant-check", "application/json",
		strings.NewReader(`{"current_participants":64}`))
	defer resp.Body.Close()
	var body participantCheckOut
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.Allowed {
		t.Fatal("at cap should not be allowed")
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
