package domain

import (
	"errors"
	"testing"
	"time"
)

func BenchmarkClassify(b *testing.B) {
	err := errors.New("dial tcp")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Classify(503, nil)
		_ = Classify(429, nil)
		_ = Classify(0, err)
		_ = Classify(0, ErrPermanentSentinel)
	}
}

func BenchmarkBackoff(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Backoff((i%10)+1, 100*time.Millisecond, 30*time.Second)
	}
}
