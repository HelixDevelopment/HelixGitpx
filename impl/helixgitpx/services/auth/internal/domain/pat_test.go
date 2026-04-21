package domain_test

import (
	"strings"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/auth/internal/domain"
)

func TestIssuePAT_FormatAndVerify(t *testing.T) {
	token, hashed, err := domain.IssuePAT()
	if err != nil {
		t.Fatalf("IssuePAT: %v", err)
	}
	if !strings.HasPrefix(token, "hpxat_") {
		t.Errorf("token missing hpxat_ prefix: %q", token)
	}
	// base62 of 24 random bytes produces 32 or 33 chars depending on the
	// leading byte; accept both widths.
	if n := len(token); n != 6+32 && n != 6+33 {
		t.Errorf("token length = %d, want 38 or 39 (prefix 6 + base62 32/33)", n)
	}
	if !domain.VerifyPAT(token, hashed) {
		t.Errorf("VerifyPAT round-trip failed")
	}
	if domain.VerifyPAT("hpxat_wrongtoken", hashed) {
		t.Errorf("VerifyPAT should reject wrong token")
	}
}

func TestVerifyPAT_RejectsNoPrefix(t *testing.T) {
	_, hashed, _ := domain.IssuePAT()
	if domain.VerifyPAT("not-a-pat", hashed) {
		t.Errorf("VerifyPAT accepted token without hpxat_ prefix")
	}
}
