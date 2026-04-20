package quota_test

import (
	"testing"
	"time"

	"github.com/helixgitpx/platform/quota"
)

func TestBucket_Allow_WithinLimit(t *testing.T) {
	b := quota.NewInMemoryBucket(5, time.Minute)
	for i := 0; i < 5; i++ {
		if !b.Allow("key") {
			t.Fatalf("allow #%d should pass", i+1)
		}
	}
	if b.Allow("key") {
		t.Errorf("6th should be denied")
	}
}

func TestBucket_Allow_DifferentKeysIndependent(t *testing.T) {
	b := quota.NewInMemoryBucket(1, time.Minute)
	if !b.Allow("a") || !b.Allow("b") {
		t.Errorf("different keys must have independent budgets")
	}
}

func TestBucket_Allow_Refill(t *testing.T) {
	b := quota.NewInMemoryBucket(1, 10*time.Millisecond)
	b.Allow("key")
	time.Sleep(15 * time.Millisecond)
	if !b.Allow("key") {
		t.Errorf("bucket should have refilled after window")
	}
}
