package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PATsPG struct{ Pool *pgxpool.Pool }

type PAT struct {
	ID        uuid.UUID
	Name      string
	Scopes    []string
	CreatedAt time.Time
	ExpiresAt *time.Time
}

func (p *PATsPG) Insert(ctx context.Context, userID, name string, hashedSecret []byte, scopes []string, expires *time.Time) (*PAT, error) {
	b, _ := json.Marshal(scopes)
	var row PAT
	err := p.Pool.QueryRow(ctx, `
		INSERT INTO auth.pats(user_id, name, hashed_secret, scopes, expires_at)
		VALUES ($1::uuid, $2, $3, $4::jsonb, $5)
		RETURNING id, name, created_at, expires_at`,
		userID, name, hashedSecret, string(b), expires,
	).Scan(&row.ID, &row.Name, &row.CreatedAt, &row.ExpiresAt)
	if err != nil {
		return nil, err
	}
	row.Scopes = scopes
	return &row, nil
}

func (p *PATsPG) Revoke(ctx context.Context, id, userID string) error {
	_, err := p.Pool.Exec(ctx,
		`UPDATE auth.pats SET revoked_at = NOW() WHERE id = $1::uuid AND user_id = $2::uuid`,
		id, userID)
	return err
}

func (p *PATsPG) List(ctx context.Context, userID string) ([]PAT, error) {
	rows, err := p.Pool.Query(ctx, `
		SELECT id, name, scopes::jsonb, created_at, expires_at
		FROM auth.pats WHERE user_id = $1::uuid AND revoked_at IS NULL`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PAT
	for rows.Next() {
		var row PAT
		var scopesJSON []byte
		if err := rows.Scan(&row.ID, &row.Name, &scopesJSON, &row.CreatedAt, &row.ExpiresAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(scopesJSON, &row.Scopes)
		out = append(out, row)
	}
	return out, nil
}
