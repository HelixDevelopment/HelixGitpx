package repo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/helixgitpx/helixgitpx/services/auth/internal/domain"
)

type SessionsPG struct{ Pool *pgxpool.Pool }

func (s *SessionsPG) Create(ctx context.Context, id uuid.UUID, userID string, expires time.Time, ua, ip string) error {
	_, err := s.Pool.Exec(ctx, `
		INSERT INTO auth.sessions(id, user_id, expires_at, user_agent, ip)
		VALUES ($1, $2::uuid, $3, $4, NULLIF($5,'')::inet)`,
		id, userID, expires, ua, ip)
	return err
}

func (s *SessionsPG) Revoke(ctx context.Context, id uuid.UUID, userID string) error {
	ct, err := s.Pool.Exec(ctx,
		`UPDATE auth.sessions SET revoked_at = NOW()
		 WHERE id = $1 AND user_id = $2::uuid AND revoked_at IS NULL`,
		id, userID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("sessions: not found or already revoked")
	}
	return nil
}

func (s *SessionsPG) List(ctx context.Context, userID string) ([]domain.Session, error) {
	rows, err := s.Pool.Query(ctx, `
		SELECT id, user_id::text, created_at, expires_at, revoked_at, user_agent, COALESCE(host(ip), '')
		FROM auth.sessions WHERE user_id = $1::uuid`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Session
	for rows.Next() {
		var s domain.Session
		if err := rows.Scan(&s.ID, &s.UserID, &s.CreatedAt, &s.ExpiresAt, &s.RevokedAt, &s.UserAgent, &s.IP); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

func (s *SessionsPG) Active(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	var row domain.Session
	err := s.Pool.QueryRow(ctx, `
		SELECT id, user_id::text, created_at, expires_at, revoked_at, user_agent, COALESCE(host(ip), '')
		FROM auth.sessions
		WHERE id = $1 AND revoked_at IS NULL AND expires_at > NOW()`, id).
		Scan(&row.ID, &row.UserID, &row.CreatedAt, &row.ExpiresAt, &row.RevokedAt, &row.UserAgent, &row.IP)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &row, err
}
