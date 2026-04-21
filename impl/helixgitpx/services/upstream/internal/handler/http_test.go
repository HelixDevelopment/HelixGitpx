package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/upstream/internal/memstore"
)

func setup(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer((&Handler{Store: memstore.New()}).Routes())
	t.Cleanup(srv.Close)
	return srv
}

func TestCreateBinding_Happy(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/upstreams", "application/json",
		strings.NewReader(`{"repo_id":"r-1","provider":"github","url":"https://github.com/o/r.git","direction":"write"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		t.Fatalf("want 201 got %d", resp.StatusCode)
	}
}

func TestCreateBinding_RejectsInvalidProvider(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/upstreams", "application/json",
		strings.NewReader(`{"repo_id":"r-1","provider":"phabricator","url":"https://x/y.git","direction":"write"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("want 400 got %d", resp.StatusCode)
	}
}

func TestList_RequiresRepoID(t *testing.T) {
	srv := setup(t)
	resp, err := http.Get(srv.URL + "/v1/upstreams")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("want 400 got %d", resp.StatusCode)
	}
}

func TestListByRepoAndDelete(t *testing.T) {
	srv := setup(t)
	create := func(provider, url string) string {
		resp, _ := http.Post(srv.URL+"/v1/upstreams", "application/json",
			strings.NewReader(`{"repo_id":"r-1","provider":"`+provider+`","url":"`+url+`","direction":"write"}`))
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
	resp, _ := http.Get(srv.URL + "/v1/upstreams?repo_id=r-1")
	defer resp.Body.Close()
	var body struct {
		Bindings []map[string]any `json:"bindings"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if len(body.Bindings) != 2 {
		t.Fatalf("want 2 bindings got %d", len(body.Bindings))
	}

	// delete
	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/v1/upstreams/"+id1, nil)
	resp2, _ := http.DefaultClient.Do(req)
	defer resp2.Body.Close()
	if resp2.StatusCode != 204 {
		t.Fatalf("want 204 got %d", resp2.StatusCode)
	}
}

func TestHealthz(t *testing.T) {
	srv := setup(t)
	resp, err := http.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal("healthz must 200")
	}
}
