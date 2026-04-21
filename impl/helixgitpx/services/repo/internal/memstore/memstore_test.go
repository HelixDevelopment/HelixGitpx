package memstore

import (
	"context"
	"errors"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/repo/internal/domain"
	"github.com/helixgitpx/helixgitpx/services/repo/internal/handler"
)

func TestCreateAssignsIDAndCreatedAt(t *testing.T) {
	s := New()
	r, err := s.Create(context.Background(), domain.Repo{OrgID: "o-1", Slug: "acme/hello"})
	if err != nil {
		t.Fatal(err)
	}
	if r.ID == "" {
		t.Fatal("ID not assigned")
	}
	if r.CreatedAt.IsZero() {
		t.Fatal("CreatedAt not populated")
	}
}

func TestGetUnknownIsErrNotFound(t *testing.T) {
	s := New()
	_, err := s.Get(context.Background(), "does-not-exist")
	if !errors.Is(err, handler.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestListFiltersByOrg(t *testing.T) {
	s := New()
	ctx := context.Background()
	_, _ = s.Create(ctx, domain.Repo{OrgID: "a", Slug: "a/x"})
	_, _ = s.Create(ctx, domain.Repo{OrgID: "a", Slug: "a/y"})
	_, _ = s.Create(ctx, domain.Repo{OrgID: "b", Slug: "b/z"})

	a, _ := s.List(ctx, "a")
	if len(a) != 2 {
		t.Fatalf("want 2 in org a, got %d", len(a))
	}
	b, _ := s.List(ctx, "b")
	if len(b) != 1 {
		t.Fatalf("want 1 in org b, got %d", len(b))
	}
	all, _ := s.List(ctx, "")
	if len(all) != 3 {
		t.Fatalf("want 3 total, got %d", len(all))
	}
}

func TestDeleteRemovesRepoAndProtections(t *testing.T) {
	s := New()
	ctx := context.Background()
	r, _ := s.Create(ctx, domain.Repo{OrgID: "o", Slug: "o/r"})
	_ = s.AddProtection(ctx, domain.Protection{RepoID: r.ID, Pattern: "refs/heads/main"})

	if err := s.Delete(ctx, r.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get(ctx, r.ID); !errors.Is(err, handler.ErrNotFound) {
		t.Fatal("repo not deleted")
	}
	if err := s.Delete(ctx, r.ID); !errors.Is(err, handler.ErrNotFound) {
		t.Fatal("delete-twice should return ErrNotFound")
	}
}

func TestProtectionsIsolatedPerRepo(t *testing.T) {
	s := New()
	ctx := context.Background()
	_ = s.AddProtection(ctx, domain.Protection{RepoID: "r-1", Pattern: "main"})
	_ = s.AddProtection(ctx, domain.Protection{RepoID: "r-1", Pattern: "release/*"})
	_ = s.AddProtection(ctx, domain.Protection{RepoID: "r-2", Pattern: "main"})

	r1, _ := s.ListProtections(ctx, "r-1")
	r2, _ := s.ListProtections(ctx, "r-2")
	if len(r1) != 2 || len(r2) != 1 {
		t.Fatalf("per-repo isolation broken: r1=%d r2=%d", len(r1), len(r2))
	}
	// Mutating the returned slice must not leak back.
	r1[0].Pattern = "tampered"
	r1b, _ := s.ListProtections(ctx, "r-1")
	if r1b[0].Pattern == "tampered" {
		t.Fatal("ListProtections returned a live reference")
	}
}
