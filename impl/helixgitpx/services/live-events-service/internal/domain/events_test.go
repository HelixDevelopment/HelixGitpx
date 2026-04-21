package domain

import (
	"errors"
	"testing"
	"time"
)

func TestResumeToken_RoundTrip(t *testing.T) {
	now := time.Unix(1700000000, 0).UTC()
	tok := ResumeToken{Offset: 42, Timestamp: now}
	s := tok.Encode()

	got, err := DecodeResumeToken(s, time.Hour, now)
	if err != nil {
		t.Fatal(err)
	}
	if got.Offset != tok.Offset || !got.Timestamp.Equal(tok.Timestamp) {
		t.Fatalf("mismatch: %+v vs %+v", got, tok)
	}
}

func TestDecodeResumeToken_Errors(t *testing.T) {
	if _, err := DecodeResumeToken("", time.Hour, time.Now()); !errors.Is(err, ErrEmptyToken) {
		t.Fatal("empty must be ErrEmptyToken")
	}
	if _, err := DecodeResumeToken("not-base64!!!", time.Hour, time.Now()); !errors.Is(err, ErrMalformedToken) {
		t.Fatal("malformed must be ErrMalformedToken")
	}
}

func TestDecodeResumeToken_Stale(t *testing.T) {
	old := time.Unix(1_000_000, 0).UTC()
	now := time.Unix(1_000_000+int64((48 * time.Hour).Seconds()), 0).UTC()
	s := ResumeToken{Offset: 1, Timestamp: old}.Encode()

	if _, err := DecodeResumeToken(s, 24*time.Hour, now); !errors.Is(err, ErrTokenStale) {
		t.Fatal("want stale error")
	}
}

func TestMatches(t *testing.T) {
	// Empty filters = match everything.
	if !Matches(nil, nil, "r1", "PUSH") {
		t.Fatal("nil filters should match")
	}
	// Filter by repo.
	if !Matches([]string{"r1", "r2"}, nil, "r1", "PUSH") {
		t.Fatal("repo in list should match")
	}
	if Matches([]string{"r1", "r2"}, nil, "r3", "PUSH") {
		t.Fatal("repo not in list must not match")
	}
	// Filter by type.
	if !Matches(nil, []string{"PUSH"}, "r1", "PUSH") {
		t.Fatal("type match expected")
	}
	if Matches(nil, []string{"PUSH"}, "r1", "CONFLICT") {
		t.Fatal("unmatched type must be filtered")
	}
}
