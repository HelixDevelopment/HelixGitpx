// Package memstore backs upstream-service with an in-memory binding store.
package memstore

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/helixgitpx/helixgitpx/services/upstream/internal/domain"
)

var ErrNotFound = errors.New("upstream: binding not found")

type Binding struct {
	ID        string
	RepoID    string
	Provider  string
	URL       string
	Direction domain.Direction
}

type Store struct {
	mu       sync.RWMutex
	bindings map[string]Binding
	byRepo   map[string][]string
	seq      atomic.Uint64
}

func New() *Store {
	return &Store{bindings: map[string]Binding{}, byRepo: map[string][]string{}}
}

func (s *Store) Create(_ context.Context, in domain.BindingInput) (Binding, error) {
	id := fmt.Sprintf("b-%08x", s.seq.Add(1))
	b := Binding{ID: id, RepoID: in.RepoID, Provider: in.Provider, URL: in.RawURL, Direction: in.Direction}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bindings[id] = b
	s.byRepo[in.RepoID] = append(s.byRepo[in.RepoID], id)
	return b, nil
}

func (s *Store) Get(_ context.Context, id string) (Binding, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.bindings[id]
	if !ok {
		return Binding{}, ErrNotFound
	}
	return b, nil
}

func (s *Store) ListByRepo(_ context.Context, repoID string) []Binding {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := s.byRepo[repoID]
	out := make([]Binding, 0, len(ids))
	for _, id := range ids {
		if b, ok := s.bindings[id]; ok {
			out = append(out, b)
		}
	}
	return out
}

func (s *Store) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.bindings[id]
	if !ok {
		return ErrNotFound
	}
	delete(s.bindings, id)
	ids := s.byRepo[b.RepoID]
	filtered := ids[:0]
	for _, x := range ids {
		if x != id {
			filtered = append(filtered, x)
		}
	}
	s.byRepo[b.RepoID] = filtered
	return nil
}
