package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewBundle_HappyPath(t *testing.T) {
	b, err := NewBundle("id1", "2.0.0", "abc123", []byte("content"), true, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if b.Hash() == "" || len(b.ETag()) < 10 {
		t.Fatal("hash/etag should be populated")
	}
	if b.Active() {
		t.Fatal("fresh bundle must not be active")
	}
}

func TestNewBundle_EmptyContent(t *testing.T) {
	_, err := NewBundle("id", "1.0", "rev", nil, false, time.Now())
	if !errors.Is(err, ErrEmptyBundle) {
		t.Fatalf("want ErrEmptyBundle, got %v", err)
	}
}

func TestNewBundle_InvalidVersion(t *testing.T) {
	for _, v := range []string{"", "not-semver", "1", "1.", ".1"} {
		if _, err := NewBundle("id", v, "rev", []byte("x"), false, time.Now()); !errors.Is(err, ErrInvalidVersion) {
			t.Errorf("%q: want ErrInvalidVersion got %v", v, err)
		}
	}
	for _, v := range []string{"1.0", "1.0.0", "1.2.3-rc.1", "1.0.0+build.5"} {
		if _, err := NewBundle("id", v, "rev", []byte("x"), false, time.Now()); err != nil {
			t.Errorf("%q: unexpected %v", v, err)
		}
	}
}

func TestCanActivate(t *testing.T) {
	signed, _ := NewBundle("s", "1.0.0", "rev", []byte("x"), true, time.Now())
	unsigned, _ := NewBundle("u", "1.0.0", "rev", []byte("x"), false, time.Now())

	if err := CanActivate(signed, true); err != nil {
		t.Fatalf("signed must activate: %v", err)
	}
	if err := CanActivate(unsigned, true); !errors.Is(err, ErrNotActiveable) {
		t.Fatal("unsigned must fail when signed required")
	}
	if err := CanActivate(unsigned, false); err != nil {
		t.Fatalf("unsigned must activate when not required: %v", err)
	}
}

func TestETagFormat(t *testing.T) {
	b, _ := NewBundle("id", "1.0.0", "rev", []byte("body"), true, time.Now())
	et := b.ETag()
	if !strings.HasPrefix(et, `"sha256-`) || !strings.HasSuffix(et, `"`) {
		t.Fatalf("etag format: %s", et)
	}
}
