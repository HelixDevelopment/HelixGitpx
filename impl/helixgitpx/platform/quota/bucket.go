package quota

import (
	"context"
	"sync"
	"time"
)

// Bucket decides if a given key is allowed to proceed.
type Bucket interface {
	Allow(key string) bool
}

// NewInMemoryBucket returns a Bucket that tracks counts in memory with a
// fixed-window algorithm. Not suitable for multi-pod deployments; use
// NewRedisBucket for production.
func NewInMemoryBucket(limit int, window time.Duration) Bucket {
	return &inMemory{limit: limit, window: window, counters: map[string]*counter{}}
}

type inMemory struct {
	mu       sync.Mutex
	limit    int
	window   time.Duration
	counters map[string]*counter
}

type counter struct {
	used     int
	resetsAt time.Time
}

func (b *inMemory) Allow(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	c, ok := b.counters[key]
	if !ok || now.After(c.resetsAt) {
		b.counters[key] = &counter{used: 1, resetsAt: now.Add(b.window)}
		return 1 <= b.limit
	}
	if c.used < b.limit {
		c.used++
		return true
	}
	return false
}

// RedisEvaler is the subset of a Redis client we need for token-bucket logic.
// Consumers pass a go-redis *redis.Client, which satisfies this via its Eval method.
type RedisEvaler interface {
	Eval(ctx context.Context, script string, keys []string, args ...any) (any, error)
}

// NewRedisBucket returns a Redis-backed bucket. Safe across pods.
func NewRedisBucket(rc RedisEvaler, prefix string, limit int, window time.Duration) Bucket {
	return &redisBucket{rc: rc, prefix: prefix, limit: limit, window: window}
}

const luaIncr = `
local c = redis.call('INCR', KEYS[1])
if c == 1 then redis.call('EXPIRE', KEYS[1], ARGV[1]) end
return c
`

type redisBucket struct {
	rc     RedisEvaler
	prefix string
	limit  int
	window time.Duration
}

func (b *redisBucket) Allow(key string) bool {
	fk := b.prefix + ":" + key
	res, err := b.rc.Eval(context.Background(), luaIncr, []string{fk}, int(b.window.Seconds()))
	if err != nil {
		return true // fail-open; M8 hardening flips to fail-closed
	}
	var n int64
	switch v := res.(type) {
	case int64:
		n = v
	case int:
		n = int64(v)
	}
	return int(n) <= b.limit
}
