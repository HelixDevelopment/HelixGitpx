package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// SessionStore is implemented by the pg repo.
type SessionStore interface {
	Create(ctx context.Context, id uuid.UUID, userID string, expires time.Time, ua, ip string) error
	Revoke(ctx context.Context, id uuid.UUID, userID string) error
	List(ctx context.Context, userID string) ([]Session, error)
	Active(ctx context.Context, id uuid.UUID) (*Session, error)
}

// Session mirrors auth.sessions.
type Session struct {
	ID        uuid.UUID
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
	RevokedAt *time.Time
	UserAgent string
	IP        string
}
