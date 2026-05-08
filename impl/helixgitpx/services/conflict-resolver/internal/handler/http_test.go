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
	srv := httptest.NewServer((&Handler{Store: NewStore()}).Routes())
	t.Cleanup(srv.Close)
	return srv
}

func TestCreateConflict_RefDivergence(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/conflicts", "application/json",
		strings.NewReader(`{"repo_id":"r-1","refs_diverge":true}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		t.Fatalf("want 201 got %d", resp.StatusCode)
	}
	var body conflictOut
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Kind != "ref_divergence" {
		t.Fatalf("want ref_divergence got %q", body.Kind)
	}
	if body.Status != "open" {
		t.Fatalf("want open got %q", body.Status)
	}
	if body.ID == "" {
		t.Fatal("expected non-empty id")
	}
}

func TestCreateConflict_NoSignal(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/conflicts", "application/json",
		strings.NewReader(`{"repo_id":"r-1"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("want 400 got %d", resp.StatusCode)
	}
	var errBody struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&errBody); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if errBody.Code != "no_conflict" {
		t.Fatalf("want code=no_conflict got %q", errBody.Code)
	}
}

func TestCreateConflict_InvalidJSON(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/conflicts", "application/json",
		strings.NewReader(`{bad`))
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
		t.Fatalf("decode: %v", err)
	}
	if errBody.Code != "invalid_json" {
		t.Fatalf("want code=invalid_json got %q", errBody.Code)
	}
}

func TestGetConflict(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/conflicts", "application/json",
		strings.NewReader(`{"repo_id":"r-1","labels_differ":true}`))
	defer resp.Body.Close()
	var created conflictOut
	_ = json.NewDecoder(resp.Body).Decode(&created)

	resp2, err := http.Get(srv.URL + "/v1/conflicts/" + created.ID)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp2.StatusCode)
	}
	var got conflictOut
	_ = json.NewDecoder(resp2.Body).Decode(&got)
	if got.Kind != "label_race" {
		t.Fatalf("want label_race got %q", got.Kind)
	}
}

func TestGetConflict_NotFound(t *testing.T) {
	srv := setup(t)
	resp, err := http.Get(srv.URL + "/v1/conflicts/nope")
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
		t.Fatalf("decode: %v", err)
	}
	if errBody.Code != "not_found" {
		t.Fatalf("want code=not_found got %q", errBody.Code)
	}
}

func TestProposeAndResolve(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/conflicts", "application/json",
		strings.NewReader(`{"repo_id":"r-1","rename_collision":true}`))
	defer resp.Body.Close()
	var created conflictOut
	_ = json.NewDecoder(resp.Body).Decode(&created)

	resp2, err := http.Post(srv.URL+"/v1/conflicts/"+created.ID+"/propose", "application/json",
		strings.NewReader(`{"rationale":"merge both"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != 200 {
		t.Fatalf("propose: want 200 got %d", resp2.StatusCode)
	}
	var proposeResult struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&proposeResult); err != nil {
		t.Fatalf("decode propose: %v", err)
	}
	if proposeResult.Status != "proposed" {
		t.Fatalf("propose: want status=proposed got %q", proposeResult.Status)
	}

	resp3, err := http.Post(srv.URL+"/v1/conflicts/"+created.ID+"/resolve", "application/json",
		strings.NewReader(`{"rationale":"merged"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp3.Body.Close()
	if resp3.StatusCode != 200 {
		t.Fatalf("resolve: want 200 got %d", resp3.StatusCode)
	}
	var resolveResult struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp3.Body).Decode(&resolveResult); err != nil {
		t.Fatalf("decode resolve: %v", err)
	}
	if resolveResult.Status != "resolved" {
		t.Fatalf("resolve: want status=resolved got %q", resolveResult.Status)
	}

	resp4, _ := http.Get(srv.URL + "/v1/conflicts/" + created.ID)
	defer resp4.Body.Close()
	var afterResolve conflictOut
	_ = json.NewDecoder(resp4.Body).Decode(&afterResolve)
	if afterResolve.Status != "resolved" {
		t.Fatalf("after resolve: want status=resolved got %q", afterResolve.Status)
	}
}

func TestPropose_EmptyRationale(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/conflicts", "application/json",
		strings.NewReader(`{"repo_id":"r-1","refs_diverge":true}`))
	defer resp.Body.Close()
	var created conflictOut
	_ = json.NewDecoder(resp.Body).Decode(&created)

	resp2, err := http.Post(srv.URL+"/v1/conflicts/"+created.ID+"/propose", "application/json",
		strings.NewReader(`{"rationale":"   "}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != 400 {
		t.Fatalf("want 400 for empty rationale got %d", resp2.StatusCode)
	}
}

func TestResolve_WithoutProposal(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/conflicts", "application/json",
		strings.NewReader(`{"repo_id":"r-1","refs_diverge":true}`))
	defer resp.Body.Close()
	var created conflictOut
	_ = json.NewDecoder(resp.Body).Decode(&created)

	resp2, err := http.Post(srv.URL+"/v1/conflicts/"+created.ID+"/resolve", "application/json",
		strings.NewReader(`{"rationale":"merged"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != 400 {
		t.Fatalf("want 400 for resolve-without-propose got %d", resp2.StatusCode)
	}
}

func TestReject(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/conflicts", "application/json",
		strings.NewReader(`{"repo_id":"r-1","meta_drift":true}`))
	defer resp.Body.Close()
	var created conflictOut
	_ = json.NewDecoder(resp.Body).Decode(&created)

	http.Post(srv.URL+"/v1/conflicts/"+created.ID+"/propose", "application/json",
		strings.NewReader(`{"rationale":"fix metadata"}`))

	resp2, err := http.Post(srv.URL+"/v1/conflicts/"+created.ID+"/reject", "application/json",
		strings.NewReader(`{"rationale":"wrong approach"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != 200 {
		t.Fatalf("reject: want 200 got %d", resp2.StatusCode)
	}
	var rejectResult struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&rejectResult); err != nil {
		t.Fatalf("decode reject: %v", err)
	}
	if rejectResult.Status != "rejected" {
		t.Fatalf("reject: want status=rejected got %q", rejectResult.Status)
	}

	resp3, _ := http.Get(srv.URL + "/v1/conflicts/" + created.ID)
	defer resp3.Body.Close()
	var afterReject conflictOut
	_ = json.NewDecoder(resp3.Body).Decode(&afterReject)
	if afterReject.Status != "rejected" {
		t.Fatalf("after reject: want status=rejected got %q", afterReject.Status)
	}
}

func TestReject_WithoutProposal(t *testing.T) {
	srv := setup(t)
	resp, _ := http.Post(srv.URL+"/v1/conflicts", "application/json",
		strings.NewReader(`{"repo_id":"r-1","refs_diverge":true}`))
	defer resp.Body.Close()
	var created conflictOut
	_ = json.NewDecoder(resp.Body).Decode(&created)

	resp2, err := http.Post(srv.URL+"/v1/conflicts/"+created.ID+"/reject", "application/json",
		strings.NewReader(`{"rationale":"nope"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != 400 {
		t.Fatalf("want 400 for reject-without-propose got %d", resp2.StatusCode)
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
