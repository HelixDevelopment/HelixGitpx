package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/orgteam/internal/domain"
)

// fakeRepo is a mock — allowed here because this file is a *unit* test.
// Integration tests must NOT use fakes; see test/integration/.
type fakeRepo struct {
	owner      string
	saved      string
	ownerErr   error
	saveErr    error
}

func (f *fakeRepo) OwnerOf(_ context.Context, _ string) (string, error) {
	return f.owner, f.ownerErr
}
func (f *fakeRepo) SetResidency(_ context.Context, _ string, r domain.Residency) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.saved = string(r)
	return nil
}

func setup(t *testing.T, repo *fakeRepo, actor string) *httptest.Server {
	t.Helper()
	h := &ResidencyHandler{
		Repo:         repo,
		ActorFromCtx: func(context.Context) string { return actor },
	}
	mux := http.NewServeMux()
	mux.Handle("POST /v1/orgs/{id}/residency", h)
	return httptest.NewServer(mux)
}

func postJSON(t *testing.T, srv *httptest.Server, path, body string) *http.Response {
	t.Helper()
	resp, err := http.Post(srv.URL+path, "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func TestResidency_OwnerChangesToEU(t *testing.T) {
	repo := &fakeRepo{owner: "alice"}
	srv := setup(t, repo, "alice")
	defer srv.Close()

	resp := postJSON(t, srv, "/v1/orgs/o-123/residency", `{"residency":"EU"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", resp.StatusCode)
	}
	var body struct{ Residency string }
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.Residency != "EU" {
		t.Fatalf("body residency=%q", body.Residency)
	}
	if repo.saved != "EU" {
		t.Fatalf("repo not updated, got %q", repo.saved)
	}
}

func TestResidency_NonOwnerIs403(t *testing.T) {
	repo := &fakeRepo{owner: "alice"}
	srv := setup(t, repo, "bob")
	defer srv.Close()

	resp := postJSON(t, srv, "/v1/orgs/o-1/residency", `{"residency":"UK"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("want 403, got %d", resp.StatusCode)
	}
}

func TestResidency_InvalidZoneIs400(t *testing.T) {
	repo := &fakeRepo{owner: "alice"}
	srv := setup(t, repo, "alice")
	defer srv.Close()

	resp := postJSON(t, srv, "/v1/orgs/o-1/residency", `{"residency":"XX"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", resp.StatusCode)
	}
}

func TestResidency_UnauthenticatedIs401(t *testing.T) {
	repo := &fakeRepo{owner: "alice"}
	srv := setup(t, repo, "") // empty actor
	defer srv.Close()

	resp := postJSON(t, srv, "/v1/orgs/o-1/residency", `{"residency":"EU"}`)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", resp.StatusCode)
	}
}
