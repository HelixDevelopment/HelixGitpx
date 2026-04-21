// Package domain encodes opa-bundle-server invariants: bundle identity,
// signature scheme, and activation rules.
package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

var (
	ErrEmptyBundle     = errors.New("bundle: content is empty")
	ErrInvalidVersion  = errors.New("bundle: version must be semver")
	ErrNotActiveable   = errors.New("bundle: cannot activate — not signed or invalid")
	ErrUnknownBundleID = errors.New("bundle: unknown id")
)

// Bundle is the in-memory handle for a .tar.gz OPA bundle. It does NOT
// keep the full content in memory; Content is passed in on construction
// and hashed immediately.
type Bundle struct {
	ID        string
	Version   string
	GitRev    string
	SHA256    [32]byte
	SizeBytes int64
	Signed    bool
	CreatedAt time.Time
	ActivatedAt time.Time
}

// ETag returns the cacheable HTTP ETag (quoted, SHA256-prefixed).
func (b Bundle) ETag() string { return `"sha256-` + hex.EncodeToString(b.SHA256[:]) + `"` }

// Active reports whether the bundle has ever been activated.
func (b Bundle) Active() bool { return !b.ActivatedAt.IsZero() }

// Hash returns the hex-encoded SHA256.
func (b Bundle) Hash() string { return hex.EncodeToString(b.SHA256[:]) }

// NewBundle constructs a Bundle from raw content + metadata.
func NewBundle(id, version, gitRev string, content []byte, signed bool, createdAt time.Time) (Bundle, error) {
	if len(content) == 0 {
		return Bundle{}, ErrEmptyBundle
	}
	if !looksSemver(version) {
		return Bundle{}, ErrInvalidVersion
	}
	sum := sha256.Sum256(content)
	return Bundle{
		ID:        id,
		Version:   version,
		GitRev:    gitRev,
		SHA256:    sum,
		SizeBytes: int64(len(content)),
		Signed:    signed,
		CreatedAt: createdAt,
	}, nil
}

// CanActivate enforces the invariant that production bundles must be signed.
func CanActivate(b Bundle, requireSigned bool) error {
	if requireSigned && !b.Signed {
		return ErrNotActiveable
	}
	return nil
}

func looksSemver(s string) bool {
	parts := strings.SplitN(s, ".", 3)
	if len(parts) < 2 {
		return false
	}
	for _, p := range parts {
		if p == "" {
			return false
		}
		core := p
		if i := strings.IndexAny(core, "-+"); i > 0 {
			core = core[:i]
		}
		for _, r := range core {
			if r < '0' || r > '9' {
				return false
			}
		}
	}
	return true
}
