// Package domain encodes git-ingress invariants applied BEFORE a push is
// accepted: ref format, ACL, and rate-limit / quota checks. These are
// enforced in addition to the Git server's own refname validation.
package domain

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrEmptyRepo      = errors.New("git: repo id is empty")
	ErrInvalidRef     = errors.New("git: invalid ref name")
	ErrForbiddenRef   = errors.New("git: ref is force-protected")
	ErrQuotaExceeded  = errors.New("git: push quota exceeded")
	ErrPushTooLarge   = errors.New("git: push exceeds per-push byte limit")
)

// RefNamePattern accepts a subset of Git's rules — enough to reject obvious
// attacks (embedded newlines, double slashes, leading dot) while letting
// legitimate branches / tags through. Strict validation is the job of the
// Git server itself; we just pre-filter.
var RefNamePattern = regexp.MustCompile(`^refs/(heads|tags|notes|pull)/[A-Za-z0-9._/\-]+$`)

// ProtectedRefs lists refs that require an extra review signal (e.g. main,
// release/*). The value of each entry is a regexp to match the ref.
var ProtectedRefs = []*regexp.Regexp{
	regexp.MustCompile(`^refs/heads/main$`),
	regexp.MustCompile(`^refs/heads/master$`),
	regexp.MustCompile(`^refs/heads/release/.*$`),
}

// ValidateRef enforces shape + protection. A protected ref still validates
// here; the ACL layer is separate.
func ValidateRef(ref string) error {
	if !RefNamePattern.MatchString(ref) {
		return ErrInvalidRef
	}
	return nil
}

// IsProtected reports whether a ref falls under a protected pattern.
func IsProtected(ref string) bool {
	for _, re := range ProtectedRefs {
		if re.MatchString(ref) {
			return true
		}
	}
	return false
}

// AllowPush decides whether a push request is allowed given the push rate
// and a per-push byte limit. The bucket caller tracks per-org usage.
type AllowPushInput struct {
	RepoID           string
	SizeBytes        int64
	PushesLastMinute int
	PushLimit        int
	MaxBytesPerPush  int64
}

func AllowPush(in AllowPushInput) error {
	if strings.TrimSpace(in.RepoID) == "" {
		return ErrEmptyRepo
	}
	if in.SizeBytes > in.MaxBytesPerPush {
		return ErrPushTooLarge
	}
	if in.PushesLastMinute >= in.PushLimit {
		return ErrQuotaExceeded
	}
	return nil
}
