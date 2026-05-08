package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/upstream/internal/memstore"
)

func setupUpstream(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer((&Handler{Store: memstore.New()}).Routes())
	t.Cleanup(srv.Close)
	return srv
}

func TestCreateBinding_Happy(t *testing.T) {
	srv := setupUpstream(t)
	resp, err := http.Post(srv.URL+"/v1/upstreams", "application/json",
		strings.NewReader(`{"repo_id":"r-1","provider":"github","url":"https://github.com/o/r.git","direction":"write"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		t.Fatalf("want 201 got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("want application/json content-type, got %q", ct)
	}
	var out bindingOut
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if out.Provider != "github" {
		t.Fatalf("want provider=github got %q", out.Provider)
	}
	if out.Direction != "write" {
		t.Fatalf("want direction=write got %q", out.Direction)
	}
	if out.RepoID != "r-1" {
		t.Fatalf("want repo_id=r-1 got %q", out.RepoID)
	}
	if out.ID == "" {
		t.Fatal("want non-empty id")
	}
}

func TestCreateBinding_RejectsInvalidProvider(t *testing.T) {
	srv := setupUpstream(t)
	resp, err := http.Post(srv.URL+"/v1/upstreams", "application/json",
		strings.NewReader(`{"repo_id":"r-1","provider":"phabricator","url":"https://x/y.git","direction":"write"}`))
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
		t.Fatalf("decode error body: %v", err)
	}
	if errBody.Code != "invalid" {
		t.Fatalf("want code=invalid got %q", errBody.Code)
	}
}

func TestCreateBinding_InvalidJSON(t *testing.T) {
	srv := setupUpstream(t)
	resp, err := http.Post(srv.URL+"/v1/upstreams", "application/json",
		strings.NewReader(`{not json`))
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

func TestList_RequiresRepoID(t *testing.T) {
	srv := setupUpstream(t)
	resp, err := http.Get(srv.URL + "/v1/upstreams")
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
	if errBody.Code != "missing_repo_id" {
		t.Fatalf("want code=missing_repo_id got %q", errBody.Code)
	}
}

func TestGetBinding_NotFound(t *testing.T) {
	srv := setupUpstream(t)
	resp, err := http.Get(srv.URL + "/v1/upstreams/does-not-exist")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatalf("want 404 got %d", resp.StatusCode)
	}
	var errBody struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&errBody); err != nil {
		t.Fatalf("decode error body: %v", err)
	}
	if errBody.Code != "not_found" {
		t.Fatalf("want code=not_found got %q", errBody.Code)
	}
}

func TestDeleteBinding_NotFound(t *testing.T) {
	srv := setupUpstream(t)
	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/v1/upstreams/does-not-exist", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatalf("want 404 got %d", resp.StatusCode)
	}
	var errBody struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&errBody); err != nil {
		t.Fatalf("decode error body: %v", err)
	}
	if errBody.Code != "not_found" {
		t.Fatalf("want code=not_found got %q", errBody.Code)
	}
}

func TestListByRepoAndDelete(t *testing.T) {
	srv := setupUpstream(t)
	create := func(provider, url string) string {
		resp, err := http.Post(srv.URL+"/v1/upstreams", "application/json",
			strings.NewReader(`{"repo_id":"r-1","provider":"`+provider+`","url":"`+url+`","direction":"write"}`))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		var out struct {
			ID string `json:"id"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&out)
		return out.ID
	}
	id1 := create("github", "https://github.com/o/r.git")
	_ = create("gitlab", "https://gitlab.com/o/r.git")

	// list
	resp, err := http.Get(srv.URL + "/v1/upstreams?repo_id=r-1")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var body struct {
		Bindings []map[string]any `json:"bindings"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if len(body.Bindings) != 2 {
		t.Fatalf("want 2 bindings got %d", len(body.Bindings))
	}

	// verify list body contains correct providers
	providers := map[string]bool{}
	for _, b := range body.Bindings {
		p, _ := b["provider"].(string)
		providers[p] = true
	}
	if !providers["github"] || !providers["gitlab"] {
		t.Fatalf("want github+gitlab in list, got %v", body.Bindings)
	}

	// get by id
	getResp, getErr := http.Get(srv.URL + "/v1/upstreams/" + id1)
	if getErr != nil {
		t.Fatal(getErr)
	}
	defer getResp.Body.Close()
	if getResp.StatusCode != 200 {
		t.Fatalf("get by id: want 200 got %d", getResp.StatusCode)
	}
	var got bindingOut
	if err := json.NewDecoder(getResp.Body).Decode(&got); err != nil {
		t.Fatalf("decode get body: %v", err)
	}
	if got.Provider != "github" {
		t.Fatalf("get by id: want provider=github got %q", got.Provider)
	}
	if got.Direction != "write" {
		t.Fatalf("get by id: want direction=write got %q", got.Direction)
	}

	// delete
	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/v1/upstreams/"+id1, nil)
	resp2, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		t.Fatal(err2)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != 204 {
		t.Fatalf("want 204 got %d", resp2.StatusCode)
	}
}

func TestHealthz_ReturnsStatusOK(t *testing.T) {
	srv := setupUpstream(t)
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
