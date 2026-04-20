package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MFAPG struct{ Pool *pgxpool.Pool }

type MFAFactor struct {
	ID     uuid.UUID
	UserID string
	Kind   string
	Secret []byte
}

func (m *MFAPG) InsertTOTP(ctx context.Context, userID, secret string) (uuid.UUID, error) {
	var id uuid.UUID
	err := m.Pool.QueryRow(ctx, `
		INSERT INTO auth.mfa_factors(user_id, kind, secret_or_pubkey)
		VALUES ($1::uuid, 'totp', $2)
		RETURNING id`, userID, []byte(secret)).Scan(&id)
	return id, err
}

func (m *MFAPG) GetTOTP(ctx context.Context, userID string) (*MFAFactor, error) {
	var f MFAFactor
	err := m.Pool.QueryRow(ctx, `
		SELECT id, user_id::text, kind::text, secret_or_pubkey
		FROM auth.mfa_factors WHERE user_id = $1::uuid AND kind = 'totp' LIMIT 1`, userID).
		Scan(&f.ID, &f.UserID, &f.Kind, &f.Secret)
	return &f, err
}
