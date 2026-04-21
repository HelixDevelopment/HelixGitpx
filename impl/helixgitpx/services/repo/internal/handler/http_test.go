package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/repo/internal/domain"
)

// fake in-memory store (unit test — mocks allowed here per Constitution §II §2).
type fakeStore struct {
	repos       map[string]domain.Repo
	protections map[string][]domain.Protection
	seq         atomic.Int32
}

func newFakeStore() *fakeStore {
	return &fakeStore{repos: map[string]domain.Repo{}, protections: map[string][]domain.Protection{}}
}

func (f *fakeStore) Create(_ context.Context, r domain.Repo) (domain.Repo, error) {
	id := "r-" + string(rune('A'+f.seq.Add(1)))
	r.ID = id
	f.repos[id] = r
	return r, nil
}

func (f *fakeStore) Get(_ context.Context, id string) (domain.Repo, error) {
	r, ok := f.repos[id]
	if !ok {
		return domain.Repo{}, ErrNotFound
	}
	return r, nil
}

func (f *fakeStore) List(_ context.Context, orgID string) ([]domain.Repo, error) {
	out := []domain.Repo{}
	for _, r := range f.repos {
		if orgID == "" || r.OrgID == orgID {
			out = append(out, r)
		}
	}
	return out, nil
}

func (f *fakeStore) Delete(_ context.Context, id string) error {
	if _, ok := f.repos[id]; !ok {
		return ErrNotFound
	}
	delete(f.repos, id)
	return nil
}

func (f *fakeStore) AddProtection(_ context.Context, p domain.Protection) error {
	f.protections[p.RepoID] = append(f.protections[p.RepoID], p)
	return nil
}

func (f *fakeStore) ListProtections(_ context.Context, repoID string) ([]domain.Protection, error) {
	return f.protections[repoID], nil
}

func setup(t *testing.T) (*httptest.Server, *fakeStore) {
	t.Helper()
	s := newFakeStore()
	srv := httptest.NewServer((&Handler{Store: s}).Routes())
	t.Cleanup(srv.Close)
	return srv, s
}

func TestCreateAndGet(t *testing.T) {
	srv, _ := setup(t)

	resp, err := http.Post(srv.URL+"/v1/repos", "application/json",
		strings.NewReader(`{"org_id":"o-1","slug":"acme/hello"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("want 201 got %d", resp.StatusCode)
	}
	var r repoOut
	_ = json.NewDecoder(resp.Body).Decode(&r)
	if r.ID == "" || r.DefaultBranch != "main" {
		t.Fatalf("unexpected %+v", r)
	}

	get, errGet := http.Get(srv.URL + "/v1/repos/" + r.ID)
	if errGet != nil { t.Fatal(errGet) }
	defer get.Body.Close()
	if get.StatusCode != http.StatusOK {
		t.Fatalf("get: want 200 got %d", get.StatusCode)
	}
}

func TestCreate_Rejects400(t *testing.T) {
	srv, _ := setup(t)
	resp, err := http.Post(srv.URL+"/v1/repos", "application/json",
		strings.NewReader(`{"slug":"no-org"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", resp.StatusCode)
	}
}

func TestList_FiltersByOrg(t *testing.T) {
	srv, s := setup(t)
	_, _ = s.Create(context.Background(), domain.Repo{OrgID: "a", Slug: "a/x"})
	_, _ = s.Create(context.Background(), domain.Repo{OrgID: "b", Slug: "b/y"})

	resp, err := http.Get(srv.URL + "/v1/repos?org_id=a")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var body struct {
		Repos []repoOut
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if len(body.Repos) != 1 || body.Repos[0].OrgID != "a" {
		t.Fatalf("filter failed: %+v", body.Repos)
	}
}

func TestGet_404(t *testing.T) {
	srv, _ := setup(t)
	resp, err := http.Get(srv.URL + "/v1/repos/does-not-exist")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("want 404 got %d", resp.StatusCode)
	}
}

func TestDelete(t *testing.T) {
	srv, s := setup(t)
	r, _ := s.Create(context.Background(), domain.Repo{OrgID: "o", Slug: "o/r"})

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/v1/repos/"+r.ID, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("want 204 got %d", resp.StatusCode)
	}
}

func TestAddProtection(t *testing.T) {
	srv, s := setup(t)
	r, _ := s.Create(context.Background(), domain.Repo{OrgID: "o", Slug: "o/r"})

	resp, err := http.Post(srv.URL+"/v1/repos/"+r.ID+"/protections", "application/json",
		strings.NewReader(`{"pattern":"refs/heads/main","require_signed":true,"required_reviewers":2}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("want 201 got %d", resp.StatusCode)
	}

	list, errList := http.Get(srv.URL + "/v1/repos/" + r.ID + "/protections")
	if errList != nil { t.Fatal(errList) }
	defer list.Body.Close()
	var body struct {
		Protections []domain.Protection
	}
	_ = json.NewDecoder(list.Body).Decode(&body)
	if len(body.Protections) != 1 || body.Protections[0].RequiredReviewers != 2 {
		t.Fatalf("protections not persisted: %+v", body.Protections)
	}
}

func TestAddProtection_RejectsNegativeReviewers(t *testing.T) {
	srv, s := setup(t)
	r, _ := s.Create(context.Background(), domain.Repo{OrgID: "o", Slug: "o/r"})
	resp, err := http.Post(srv.URL+"/v1/repos/"+r.ID+"/protections", "application/json",
		strings.NewReader(`{"pattern":"refs/heads/*","required_reviewers":-1}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", resp.StatusCode)
	}
}

func TestHealthz(t *testing.T) {
	srv, _ := setup(t)
	resp, err := http.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatal("healthz must be 200")
	}
}
