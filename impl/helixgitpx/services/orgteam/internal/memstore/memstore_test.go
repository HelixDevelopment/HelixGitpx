package memstore

import (
	"context"
	"errors"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/orgteam/internal/domain"
)

func TestCreateOrg_AssignsIDAndIndexesSlug(t *testing.T) {
	s := New()
	o, err := s.CreateOrg(context.Background(), "acme", "Acme Inc.", "alice", domain.ResidencyEU)
	if err != nil {
		t.Fatal(err)
	}
	if o.ID == "" {
		t.Fatal("ID not assigned")
	}
	if o.Slug != "acme" {
		t.Fatalf("slug mismatch: %q", o.Slug)
	}
}

func TestCreateOrg_RejectsDuplicateSlug(t *testing.T) {
	s := New()
	_, _ = s.CreateOrg(context.Background(), "acme", "Acme Inc.", "alice", domain.ResidencyEU)
	_, err := s.CreateOrg(context.Background(), "acme", "Other Acme", "bob", domain.ResidencyEU)
	if !errors.Is(err, ErrDuplicateSlug) {
		t.Fatalf("want ErrDuplicateSlug got %v", err)
	}
}

func TestOwnerOf(t *testing.T) {
	s := New()
	o, _ := s.CreateOrg(context.Background(), "acme", "Acme", "alice", domain.ResidencyEU)
	owner, err := s.OwnerOf(context.Background(), o.ID)
	if err != nil || owner != "alice" {
		t.Fatalf("owner drift: %v %q", err, owner)
	}
	if _, err := s.OwnerOf(context.Background(), "unknown"); !errors.Is(err, ErrOrgNotFound) {
		t.Fatalf("unknown org should return ErrOrgNotFound, got %v", err)
	}
}

func TestSetResidency(t *testing.T) {
	s := New()
	o, _ := s.CreateOrg(context.Background(), "acme", "Acme", "alice", domain.ResidencyEU)
	if err := s.SetResidency(context.Background(), o.ID, domain.ResidencyUK); err != nil {
		t.Fatal(err)
	}
	_, rez, _ := s.GetOrg(context.Background(), o.ID)
	if rez != domain.ResidencyUK {
		t.Fatalf("want UK got %q", rez)
	}
	if err := s.SetResidency(context.Background(), "unknown", domain.ResidencyUS); !errors.Is(err, ErrOrgNotFound) {
		t.Fatalf("unknown should return ErrOrgNotFound, got %v", err)
	}
}

func TestListOrgs(t *testing.T) {
	s := New()
	ctx := context.Background()
	_, _ = s.CreateOrg(ctx, "acme", "Acme", "alice", domain.ResidencyEU)
	_, _ = s.CreateOrg(ctx, "contoso", "Contoso", "bob", domain.ResidencyUK)
	_, _ = s.CreateOrg(ctx, "fabrikam", "Fabrikam", "carol", domain.ResidencyUS)
	if got := len(s.ListOrgs(ctx)); got != 3 {
		t.Fatalf("want 3 orgs got %d", got)
	}
}
