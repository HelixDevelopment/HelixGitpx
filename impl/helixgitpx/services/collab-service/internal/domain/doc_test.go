package domain

import (
	"errors"
	"testing"
)

func TestValidateOpenDoc(t *testing.T) {
	cases := map[string]struct {
		docID, actor string
		want         error
	}{
		"happy":      {"d1", "alice", nil},
		"empty doc":  {"  ", "alice", ErrEmptyDocID},
		"empty actor": {"d1", "", ErrEmptyActor},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := ValidateOpenDoc(c.docID, c.actor)
			if !errors.Is(got, c.want) {
				t.Fatalf("want %v got %v", c.want, got)
			}
		})
	}
}

func TestSnapshotSizeAllowed(t *testing.T) {
	lim := Limits{MaxSnapshotBytes: 4}
	if err := SnapshotSizeAllowed(lim, []byte("abcd")); err != nil {
		t.Fatal(err)
	}
	if err := SnapshotSizeAllowed(lim, []byte("abcde")); !errors.Is(err, ErrDocumentTooLarge) {
		t.Fatalf("expected ErrDocumentTooLarge, got %v", err)
	}
}

func TestAddParticipantAllowed(t *testing.T) {
	lim := Limits{MaxParticipants: 2}
	if err := AddParticipantAllowed(lim, 1); err != nil {
		t.Fatal(err)
	}
	if err := AddParticipantAllowed(lim, 2); !errors.Is(err, ErrTooManyParticipants) {
		t.Fatalf("expected rejection at cap")
	}
}

func TestDefaultLimits(t *testing.T) {
	d := DefaultLimits()
	if d.MaxSnapshotBytes == 0 || d.MaxParticipants == 0 || d.IdleTimeout == 0 {
		t.Fatalf("defaults must be non-zero: %+v", d)
	}
}
