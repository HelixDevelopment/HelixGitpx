package domain

import (
	"errors"
	"testing"
	"time"
)

func TestClassify(t *testing.T) {
	if Classify(429, nil) != KindRateLimit {
		t.Fatal("429 must be rate-limit")
	}
	if Classify(503, nil) != KindTransient {
		t.Fatal("5xx must be transient")
	}
	if Classify(401, nil) != KindAuthFailed {
		t.Fatal("401 must be auth")
	}
	if Classify(404, nil) != KindClientError {
		t.Fatal("404 must be client")
	}
	if Classify(0, errors.New("dial tcp")) != KindTransient {
		t.Fatal("network error must be transient")
	}
	if Classify(0, ErrPermanentSentinel) != KindPermanent {
		t.Fatal("permanent sentinel must route to KindPermanent")
	}
}

func TestBackoff(t *testing.T) {
	if got := Backoff(1, time.Second, time.Minute); got != time.Second {
		t.Fatalf("attempt 1: want 1s got %v", got)
	}
	if got := Backoff(4, time.Second, time.Minute); got != 8*time.Second {
		t.Fatalf("attempt 4: want 8s got %v", got)
	}
	if got := Backoff(10, time.Second, 30*time.Second); got != 30*time.Second {
		t.Fatalf("cap not applied: %v", got)
	}
	if got := Backoff(0, time.Second, time.Minute); got != 0 {
		t.Fatalf("non-positive attempt: want 0 got %v", got)
	}
}

func TestShouldRetry(t *testing.T) {
	if !ShouldRetry(KindTransient, 1, 5) {
		t.Fatal("transient below cap must retry")
	}
	if ShouldRetry(KindAuthFailed, 1, 5) {
		t.Fatal("auth-failed must not retry")
	}
	if ShouldRetry(KindTransient, 5, 5) {
		t.Fatal("at cap must not retry")
	}
}

func TestGoesToDLQ(t *testing.T) {
	if !GoesToDLQ(KindPermanent, 1, 5) {
		t.Fatal("permanent must DLQ immediately")
	}
	if !GoesToDLQ(KindClientError, 1, 5) {
		t.Fatal("client error must DLQ immediately")
	}
	if GoesToDLQ(KindTransient, 1, 5) {
		t.Fatal("transient below cap must NOT DLQ")
	}
	if !GoesToDLQ(KindTransient, 5, 5) {
		t.Fatal("transient at cap must DLQ")
	}
}
