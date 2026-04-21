// Package domain holds the collab-service's CRDT-adjacent invariants.
// The Automerge-go library handles the actual CRDT math; this package
// encodes the rules HelixGitpx applies on top: access, rate-limiting,
// and snapshot eviction.
package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrEmptyDocID         = errors.New("collab: doc id is empty")
	ErrEmptyActor         = errors.New("collab: actor id is empty")
	ErrDocumentTooLarge   = errors.New("collab: document exceeds size limit")
	ErrTooManyParticipants = errors.New("collab: too many concurrent participants")
)

// Limits bound what the service will accept per-document. Deployable via Helm.
type Limits struct {
	MaxSnapshotBytes int
	MaxParticipants  int
	IdleTimeout      time.Duration
}

// DefaultLimits are the values applied when the operator hasn't overridden.
func DefaultLimits() Limits {
	return Limits{
		MaxSnapshotBytes: 8 * 1024 * 1024,
		MaxParticipants:  64,
		IdleTimeout:      30 * time.Minute,
	}
}

// ValidateOpenDoc checks the invariants that can be enforced without
// touching Postgres or Automerge.
func ValidateOpenDoc(docID, actorID string) error {
	if strings.TrimSpace(docID) == "" {
		return ErrEmptyDocID
	}
	if strings.TrimSpace(actorID) == "" {
		return ErrEmptyActor
	}
	return nil
}

// SnapshotSizeAllowed reports whether a candidate snapshot is within limits.
func SnapshotSizeAllowed(limits Limits, snapshot []byte) error {
	if len(snapshot) > limits.MaxSnapshotBytes {
		return ErrDocumentTooLarge
	}
	return nil
}

// AddParticipantAllowed reports whether a new participant may join.
func AddParticipantAllowed(limits Limits, current int) error {
	if current >= limits.MaxParticipants {
		return ErrTooManyParticipants
	}
	return nil
}
