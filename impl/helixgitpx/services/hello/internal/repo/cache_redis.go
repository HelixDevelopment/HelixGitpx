package repo

import (
	"context"
	"time"

	hr "github.com/helixgitpx/platform/redis"
)

// CacheRedis implements domain.Cache using HelixGitpx's Redis wrapper.
type CacheRedis struct {
	Client *hr.Client
	TTL    time.Duration
}

func (c *CacheRedis) SetLast(ctx context.Context, name, greeting string) error {
	key := c.Client.Key("last", name)
	ttl := c.TTL
	if ttl == 0 {
		ttl = 10 * time.Minute
	}
	return c.Client.Set(ctx, key, greeting, ttl).Err()
}
