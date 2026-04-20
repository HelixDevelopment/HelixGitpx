package domain_test

import (
	"testing"

	"github.com/helixgitpx/helixgitpx/services/repo/internal/domain"
)

func TestMatchesPattern(t *testing.T) {
	cases := []struct {
		pattern, ref string
		want         bool
	}{
		{"main", "refs/heads/main", true},
		{"main", "refs/heads/feature", false},
		{"release/*", "refs/heads/release/v1", true},
		{"release/*", "refs/heads/main", false},
		{"*", "refs/heads/anything", true},
	}
	for _, c := range cases {
		if got := domain.MatchesPattern(c.pattern, c.ref); got != c.want {
			t.Errorf("MatchesPattern(%q, %q) = %v, want %v", c.pattern, c.ref, got, c.want)
		}
	}
}
