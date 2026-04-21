// Package domain encodes the live-events-service invariants: resume-token
// format, subscriber filter matching, and fan-out fairness.
package domain

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"time"
)

var (
	ErrEmptyToken   = errors.New("events: resume token is empty")
	ErrMalformedToken = errors.New("events: resume token malformed")
	ErrTokenStale   = errors.New("events: resume token older than retention window")
)

// ResumeToken is an opaque bearer for the event stream position.
// Wire format: base64( 8 bytes kafka offset | 8 bytes unix-seconds ).
type ResumeToken struct {
	Offset    int64
	Timestamp time.Time
}

// Encode serialises into the on-wire format.
func (r ResumeToken) Encode() string {
	var buf [16]byte
	binary.BigEndian.PutUint64(buf[:8], uint64(r.Offset))
	binary.BigEndian.PutUint64(buf[8:], uint64(r.Timestamp.Unix()))
	return base64.RawURLEncoding.EncodeToString(buf[:])
}

// DecodeResumeToken parses a wire-format token and rejects it if older than
// `retention`. An empty string is an explicit error (use zero-value
// ResumeToken for "start from beginning").
func DecodeResumeToken(s string, retention time.Duration, now time.Time) (ResumeToken, error) {
	if s == "" {
		return ResumeToken{}, ErrEmptyToken
	}
	raw, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil || len(raw) != 16 {
		return ResumeToken{}, ErrMalformedToken
	}
	t := ResumeToken{
		Offset:    int64(binary.BigEndian.Uint64(raw[:8])),
		Timestamp: time.Unix(int64(binary.BigEndian.Uint64(raw[8:])), 0).UTC(),
	}
	if retention > 0 && now.Sub(t.Timestamp) > retention {
		return t, ErrTokenStale
	}
	return t, nil
}

// Matches reports whether an event with the given repo-id and type should be
// delivered to a subscriber with the given filter. An empty repo list means
// "any repo"; an empty type list means "any type".
func Matches(subRepos, subTypes []string, repoID, eventType string) bool {
	if len(subRepos) > 0 && !contains(subRepos, repoID) {
		return false
	}
	if len(subTypes) > 0 && !contains(subTypes, eventType) {
		return false
	}
	return true
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
