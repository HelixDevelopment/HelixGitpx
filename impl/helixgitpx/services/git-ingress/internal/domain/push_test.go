package domain

import (
	"errors"
	"testing"
)

func TestValidateRef(t *testing.T) {
	good := []string{"refs/heads/main", "refs/tags/v1.2.3", "refs/pull/42/head"}
	for _, r := range good {
		if err := ValidateRef(r); err != nil {
			t.Errorf("%q: unexpected %v", r, err)
		}
	}
	bad := []string{"", "main", "refs/heads/", "refs/heads/with space", "refs/random/foo"}
	for _, r := range bad {
		if err := ValidateRef(r); !errors.Is(err, ErrInvalidRef) {
			t.Errorf("%q: want ErrInvalidRef got %v", r, err)
		}
	}
}

func TestIsProtected(t *testing.T) {
	if !IsProtected("refs/heads/main") {
		t.Fatal("main must be protected")
	}
	if !IsProtected("refs/heads/release/2026-04") {
		t.Fatal("release/* must be protected")
	}
	if IsProtected("refs/heads/feat/x") {
		t.Fatal("feat branches must not be protected")
	}
}

func TestAllowPush(t *testing.T) {
	base := AllowPushInput{
		RepoID:          "repo-1",
		SizeBytes:       1000,
		PushesLastMinute: 5,
		PushLimit:       10,
		MaxBytesPerPush: 100_000,
	}
	if err := AllowPush(base); err != nil {
		t.Fatalf("unexpected %v", err)
	}

	tooBig := base
	tooBig.SizeBytes = base.MaxBytesPerPush + 1
	if err := AllowPush(tooBig); !errors.Is(err, ErrPushTooLarge) {
		t.Fatal("want ErrPushTooLarge")
	}

	overQuota := base
	overQuota.PushesLastMinute = base.PushLimit
	if err := AllowPush(overQuota); !errors.Is(err, ErrQuotaExceeded) {
		t.Fatal("want ErrQuotaExceeded")
	}

	emptyRepo := base
	emptyRepo.RepoID = ""
	if err := AllowPush(emptyRepo); !errors.Is(err, ErrEmptyRepo) {
		t.Fatal("want ErrEmptyRepo")
	}
}
