// Package store holds OPA bundle bytes in memory, keyed by ID, with a
// pointer to the currently-active bundle. Real deployments back this with
// MinIO; this in-memory impl serves local dev + tests.
package store

import (
	"errors"
	"sync"

	"github.com/helixgitpx/helixgitpx/services/opa-bundle-server/internal/domain"
)

var ErrInactiveRequested = errors.New("store: no active bundle")

type Store struct {
	mu      sync.RWMutex
	byID    map[string]record
	active  string
}

type record struct {
	meta    domain.Bundle
	content []byte
}

func New() *Store {
	return &Store{byID: map[string]record{}}
}

// Put stores a bundle and its content. Returns the updated metadata.
func (s *Store) Put(meta domain.Bundle, content []byte) domain.Bundle {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byID[meta.ID] = record{meta: meta, content: append([]byte(nil), content...)}
	return meta
}

// Activate marks a bundle as active and records the activation time.
func (s *Store) Activate(id string, activatedMeta domain.Bundle) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.byID[id]
	if !ok {
		return domain.ErrUnknownBundleID
	}
	r.meta = activatedMeta
	s.byID[id] = r
	s.active = id
	return nil
}

// Active returns the currently-active bundle's meta + content, or an error.
func (s *Store) Active() (domain.Bundle, []byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.active == "" {
		return domain.Bundle{}, nil, ErrInactiveRequested
	}
	r := s.byID[s.active]
	return r.meta, r.content, nil
}

// Get returns meta + content for a given id.
func (s *Store) Get(id string) (domain.Bundle, []byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.byID[id]
	if !ok {
		return domain.Bundle{}, nil, domain.ErrUnknownBundleID
	}
	return r.meta, r.content, nil
}

// List returns every bundle's metadata (newest-first by creation time).
func (s *Store) List() []domain.Bundle {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Bundle, 0, len(s.byID))
	for _, r := range s.byID {
		out = append(out, r.meta)
	}
	return out
}
