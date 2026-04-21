// Package memstore backs the orgteam handlers with an in-memory Store.
// Swappable for a Postgres implementation when the org.organizations
// schema bootstrap is wired.
package memstore

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/helixgitpx/helixgitpx/services/orgteam/internal/domain"
	"github.com/helixgitpx/helixgitpx/services/orgteam/internal/handler"
)

var ErrDuplicateSlug = errors.New("orgteam: slug already exists")
var ErrOrgNotFound = errors.New("orgteam: org not found")

type orgRecord struct {
	domain.Org
	OwnerID   string
	Residency domain.Residency
}

type Store struct {
	mu    sync.RWMutex
	orgs  map[string]orgRecord
	slugs map[string]string
	seq   atomic.Uint64
}

func New() *Store {
	return &Store{orgs: map[string]orgRecord{}, slugs: map[string]string{}}
}

// CreateOrg registers a new org with the given owner.
func (s *Store) CreateOrg(_ context.Context, slug, name, ownerID string, residency domain.Residency) (domain.Org, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.slugs[slug]; exists {
		return domain.Org{}, ErrDuplicateSlug
	}
	id := fmt.Sprintf("o-%08x", s.seq.Add(1))
	rec := orgRecord{
		Org:       domain.Org{ID: id, Slug: slug, Name: name},
		OwnerID:   ownerID,
		Residency: residency,
	}
	s.orgs[id] = rec
	s.slugs[slug] = id
	return rec.Org, nil
}

// OwnerOf implements handler.OrgRepository.
func (s *Store) OwnerOf(_ context.Context, orgID string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rec, ok := s.orgs[orgID]
	if !ok {
		return "", ErrOrgNotFound
	}
	return rec.OwnerID, nil
}

// SetResidency implements handler.OrgRepository.
func (s *Store) SetResidency(_ context.Context, orgID string, residency domain.Residency) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.orgs[orgID]
	if !ok {
		return ErrOrgNotFound
	}
	rec.Residency = residency
	s.orgs[orgID] = rec
	return nil
}

// GetOrg returns the org record (plus residency).
func (s *Store) GetOrg(_ context.Context, orgID string) (domain.Org, domain.Residency, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rec, ok := s.orgs[orgID]
	if !ok {
		return domain.Org{}, "", ErrOrgNotFound
	}
	return rec.Org, rec.Residency, nil
}

// ListOrgs returns every org.
func (s *Store) ListOrgs(_ context.Context) []domain.Org {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Org, 0, len(s.orgs))
	for _, rec := range s.orgs {
		out = append(out, rec.Org)
	}
	return out
}

// Ensure we satisfy the handler.OrgRepository interface.
var _ handler.OrgRepository = (*Store)(nil)
