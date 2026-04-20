// Package domain (tokens) orchestrates JWT issuance + refresh rotation.
package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/helixgitpx/platform/auth"
)

// Tokens bundles an access token + rotating refresh id.
type Tokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Duration
}

// TokensIssuer mints token pairs. Refresh tokens are opaque session ids.
type TokensIssuer struct {
	Signer     *auth.Signer
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

// Issue mints a new token pair for user uid, recording a session row via persist.
func (t *TokensIssuer) Issue(_ context.Context, uid string, persist func(sessionID uuid.UUID, expires time.Time) error) (*Tokens, error) {
	accessTok, err := t.Signer.Issue(auth.Claims{Subject: uid, TTL: t.AccessTTL})
	if err != nil {
		return nil, err
	}
	sessionID := uuid.New()
	refresh := sessionID.String()
	expires := time.Now().Add(t.RefreshTTL)
	if err := persist(sessionID, expires); err != nil {
		return nil, fmt.Errorf("tokens: persist session: %w", err)
	}
	return &Tokens{
		AccessToken:  accessTok,
		RefreshToken: refresh,
		ExpiresIn:    t.AccessTTL,
	}, nil
}
