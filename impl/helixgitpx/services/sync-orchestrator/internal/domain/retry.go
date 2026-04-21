// Package domain encodes sync-orchestrator invariants: attempt classification,
// retry back-off, and DLQ admission.
package domain

import (
	"errors"
	"math"
	"time"
)

type ErrorKind int

const (
	KindUnknown ErrorKind = iota
	KindTransient          // 5xx, timeout, network hiccup
	KindRateLimit          // 429
	KindAuthFailed         // 401/403
	KindClientError        // other 4xx — not retriable
	KindPermanent          // explicit provider "will never succeed"
)

// Classify translates an HTTP status code (or a sentinel error) into an ErrorKind.
func Classify(httpStatus int, err error) ErrorKind {
	if errors.Is(err, ErrPermanentSentinel) {
		return KindPermanent
	}
	switch {
	case httpStatus == 429:
		return KindRateLimit
	case httpStatus == 401 || httpStatus == 403:
		return KindAuthFailed
	case httpStatus >= 500:
		return KindTransient
	case httpStatus >= 400:
		return KindClientError
	}
	if err != nil {
		return KindTransient
	}
	return KindUnknown
}

// ErrPermanentSentinel marks errors that should never retry.
var ErrPermanentSentinel = errors.New("permanent")

// Backoff returns the next delay given attempt number (starting at 1).
// Exponential: 2^(attempt-1) * base, capped at max.
func Backoff(attempt int, base, max time.Duration) time.Duration {
	if attempt <= 0 {
		return 0
	}
	d := time.Duration(math.Pow(2, float64(attempt-1))) * base
	if d > max {
		return max
	}
	return d
}

// ShouldRetry reports whether another attempt should be made.
func ShouldRetry(kind ErrorKind, attempt, maxAttempts int) bool {
	if attempt >= maxAttempts {
		return false
	}
	switch kind {
	case KindTransient, KindRateLimit:
		return true
	default:
		return false
	}
}

// GoesToDLQ reports whether a failed job should land in the dead-letter queue.
func GoesToDLQ(kind ErrorKind, attempt, maxAttempts int) bool {
	switch kind {
	case KindPermanent, KindClientError, KindAuthFailed:
		return true
	case KindTransient, KindRateLimit:
		return attempt >= maxAttempts
	default:
		return attempt >= maxAttempts
	}
}
