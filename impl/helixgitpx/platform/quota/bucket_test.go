package quota_test

import (
	"sync"
	"sync/atomic"
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

func TestBucket_Allow_ZeroLimit_DisallowsAll(t *testing.T) {
	b := quota.NewInMemoryBucket(0, time.Minute)
	if b.Allow("key") {
		t.Error("zero-limit bucket should deny first request")
	}
}

func TestBucket_Allow_ConcurrentNoRace(t *testing.T) {
	b := quota.NewInMemoryBucket(1000, time.Minute)
	var allowed, denied atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 2000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if b.Allow("shared") {
				allowed.Add(1)
			} else {
				denied.Add(1)
			}
		}()
	}
	wg.Wait()
	a := allowed.Load()
	d := denied.Load()
	if a != 1000 {
		t.Errorf("allowed = %d, want 1000", a)
	}
	if d != 1000 {
		t.Errorf("denied = %d, want 1000", d)
	}
}

func TestBucket_Allow_RefillRestoresCapacity(t *testing.T) {
	b := quota.NewInMemoryBucket(3, 20*time.Millisecond)
	if !b.Allow("k") || !b.Allow("k") || !b.Allow("k") {
		t.Fatal("first 3 should pass")
	}
	if b.Allow("k") {
		t.Fatal("4th should be denied before refill")
	}
	time.Sleep(25 * time.Millisecond)
	if !b.Allow("k") {
		t.Error("should pass after window reset")
	}
}

func TestBucket_Allow_MultipleKeysExhaustIndependently(t *testing.T) {
	b := quota.NewInMemoryBucket(1, time.Minute)
	if !b.Allow("alpha") {
		t.Error("alpha first call should pass")
	}
	if b.Allow("alpha") {
		t.Error("alpha second call should be denied")
	}
	if !b.Allow("beta") {
		t.Error("beta should have its own budget")
	}
}
