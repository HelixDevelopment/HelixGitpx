package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersPG struct{ Pool *pgxpool.Pool }

type User struct {
	ID          uuid.UUID
	Subject     string
	Email       string
	DisplayName string
}

func (u *UsersPG) UpsertBySubject(ctx context.Context, subject, email, displayName string) (*User, error) {
	var row User
	err := u.Pool.QueryRow(ctx, `
		INSERT INTO auth.users(subject, email, display_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (subject) DO UPDATE
		  SET email = EXCLUDED.email, display_name = EXCLUDED.display_name
		RETURNING id, subject, email, display_name`,
		subject, email, displayName,
	).Scan(&row.ID, &row.Subject, &row.Email, &row.DisplayName)
	if err != nil {
		return nil, fmt.Errorf("users: upsert: %w", err)
	}
	return &row, nil
}

func (u *UsersPG) GetBySubject(ctx context.Context, subject string) (*User, error) {
	var row User
	err := u.Pool.QueryRow(ctx, `SELECT id, subject, email, display_name FROM auth.users WHERE subject = $1`, subject).
		Scan(&row.ID, &row.Subject, &row.Email, &row.DisplayName)
	if err != nil {
		return nil, err
	}
	return &row, nil
}
