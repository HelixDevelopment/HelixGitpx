package memstore

import (
	"context"
	"errors"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/upstream/internal/domain"
)

func mk(repo, provider, url string) domain.BindingInput {
	return domain.BindingInput{
		RepoID:    repo,
		Provider:  provider,
		RawURL:    url,
		Direction: domain.DirectionWrite,
	}
}

func TestCreateAndGet(t *testing.T) {
	s := New()
	b, err := s.Create(context.Background(), mk("r-1", "github", "https://github.com/o/r.git"))
	if err != nil {
		t.Fatal(err)
	}
	if b.ID == "" {
		t.Fatal("ID not assigned")
	}

	got, err := s.Get(context.Background(), b.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Provider != "github" {
		t.Fatalf("provider drift: %s", got.Provider)
	}
}

func TestGetUnknownIsErrNotFound(t *testing.T) {
	s := New()
	_, err := s.Get(context.Background(), "nope")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound got %v", err)
	}
}

func TestListByRepo(t *testing.T) {
	s := New()
	ctx := context.Background()
	_, _ = s.Create(ctx, mk("r-1", "github", "https://github.com/o/r.git"))
	_, _ = s.Create(ctx, mk("r-1", "gitlab", "https://gitlab.com/o/r.git"))
	_, _ = s.Create(ctx, mk("r-2", "gitea", "https://gitea.io/o/r.git"))

	r1 := s.ListByRepo(ctx, "r-1")
	if len(r1) != 2 {
		t.Fatalf("want 2 for r-1, got %d", len(r1))
	}
	r2 := s.ListByRepo(ctx, "r-2")
	if len(r2) != 1 {
		t.Fatalf("want 1 for r-2, got %d", len(r2))
	}
	if got := s.ListByRepo(ctx, "unknown"); len(got) != 0 {
		t.Fatal("unknown repo should return empty")
	}
}

func TestDelete(t *testing.T) {
	s := New()
	ctx := context.Background()
	a, _ := s.Create(ctx, mk("r-1", "github", "https://github.com/o/r.git"))
	b, _ := s.Create(ctx, mk("r-1", "gitlab", "https://gitlab.com/o/r.git"))

	if err := s.Delete(ctx, a.ID); err != nil {
		t.Fatal(err)
	}
	remaining := s.ListByRepo(ctx, "r-1")
	if len(remaining) != 1 || remaining[0].ID != b.ID {
		t.Fatalf("delete broke repo index: %v", remaining)
	}
	if err := s.Delete(ctx, a.ID); !errors.Is(err, ErrNotFound) {
		t.Fatal("delete-twice should return ErrNotFound")
	}
}
