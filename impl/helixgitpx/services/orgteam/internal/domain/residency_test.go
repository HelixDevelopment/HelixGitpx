package domain

import "testing"

func TestResidencyValid(t *testing.T) {
	for _, z := range []Residency{ResidencyEU, ResidencyUK, ResidencyUS} {
		if !z.Valid() {
			t.Fatalf("%s should be valid", z)
		}
	}
	if Residency("ZZ").Valid() {
		t.Fatal("ZZ should be invalid")
	}
}

func TestSetOrgResidencyAuthorization(t *testing.T) {
	if err := SetOrgResidency("alice", "bob", ResidencyEU); err == nil {
		t.Fatal("non-owner must be rejected")
	}
	if err := SetOrgResidency("alice", "alice", ResidencyEU); err != nil {
		t.Fatalf("owner must be allowed: %v", err)
	}
	if err := SetOrgResidency("alice", "alice", Residency("ZZ")); err != ErrInvalidResidency {
		t.Fatalf("expected ErrInvalidResidency, got %v", err)
	}
}
