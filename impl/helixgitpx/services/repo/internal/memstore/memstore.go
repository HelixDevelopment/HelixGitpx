// Package memstore provides an in-memory implementation of handler.Store.
// It is the default backing store for local dev and `make dev`. The
// Postgres-backed store lands in a sibling package and is selected at
// boot-time via env vars; for now, the in-memory store is authoritative
// for the GA tag.
package memstore

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/helixgitpx/helixgitpx/services/repo/internal/domain"
	"github.com/helixgitpx/helixgitpx/services/repo/internal/handler"
)

type Store struct {
	mu          sync.RWMutex
	repos       map[string]domain.Repo
	protections map[string][]domain.Protection
	seq         atomic.Uint64
}

func New() *Store {
	return &Store{
		repos:       map[string]domain.Repo{},
		protections: map[string][]domain.Protection{},
	}
}

func (s *Store) Create(_ context.Context, r domain.Repo) (domain.Repo, error) {
	id := fmt.Sprintf("r-%08x", s.seq.Add(1))
	r.ID = id
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now().UTC()
	}
	s.mu.Lock()
	s.repos[id] = r
	s.mu.Unlock()
	return r, nil
}

func (s *Store) Get(_ context.Context, id string) (domain.Repo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.repos[id]
	if !ok {
		return domain.Repo{}, handler.ErrNotFound
	}
	return r, nil
}

func (s *Store) List(_ context.Context, orgID string) ([]domain.Repo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Repo, 0, len(s.repos))
	for _, r := range s.repos {
		if orgID == "" || r.OrgID == orgID {
			out = append(out, r)
		}
	}
	return out, nil
}

func (s *Store) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.repos[id]; !ok {
		return handler.ErrNotFound
	}
	delete(s.repos, id)
	delete(s.protections, id)
	return nil
}

func (s *Store) AddProtection(_ context.Context, p domain.Protection) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.protections[p.RepoID] = append(s.protections[p.RepoID], p)
	return nil
}

func (s *Store) ListProtections(_ context.Context, repoID string) ([]domain.Protection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]domain.Protection(nil), s.protections[repoID]...), nil
}
