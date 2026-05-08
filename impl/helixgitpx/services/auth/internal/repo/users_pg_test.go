package repo

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func openAuthPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	d := os.Getenv("TEST_PG_DSN")
	if d == "" {
		d = "postgres://helix:helix@localhost:15432/helixgitpx?sslmode=disable"
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, d)
	if err != nil {
		t.Fatalf("open pool: %v", err)
	}
	t.Cleanup(pool.Close)
	if err := pool.Ping(ctx); err != nil {
		t.Skipf("SKIP-OK: #integration — Postgres not available: %v", err)
	}
	return pool
}

func TestUsersPG_UpsertBySubject_Insert(t *testing.T) {
	pool := openAuthPool(t)
	ctx := context.Background()
	u := &UsersPG{Pool: pool}

	subject := "test|" + uuid.New().String()
	user, err := u.UpsertBySubject(ctx, subject, "test@example.com", "Test User")
	if err != nil {
		t.Fatalf("UpsertBySubject: %v", err)
	}
	if user.ID == uuid.Nil {
		t.Error("ID should not be nil UUID")
	}
	if user.Subject != subject {
		t.Errorf("Subject = %q, want %q", user.Subject, subject)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Email = %q", user.Email)
	}
	if user.DisplayName != "Test User" {
		t.Errorf("DisplayName = %q", user.DisplayName)
	}

	_, _ = pool.Exec(ctx, "DELETE FROM auth.users WHERE id = $1", user.ID)
}

func TestUsersPG_UpsertBySubject_UpdateExisting(t *testing.T) {
	pool := openAuthPool(t)
	ctx := context.Background()
	u := &UsersPG{Pool: pool}

	subject := "test|" + uuid.New().String()
	user1, _ := u.UpsertBySubject(ctx, subject, "old@example.com", "Old Name")
	user2, err := u.UpsertBySubject(ctx, subject, "new@example.com", "New Name")
	if err != nil {
		t.Fatalf("UpsertBySubject update: %v", err)
	}
	if user2.ID != user1.ID {
		t.Errorf("ID changed on upsert: %s -> %s", user1.ID, user2.ID)
	}
	if user2.Email != "new@example.com" {
		t.Errorf("Email = %q, want new@example.com", user2.Email)
	}
	if user2.DisplayName != "New Name" {
		t.Errorf("DisplayName = %q, want New Name", user2.DisplayName)
	}

	_, _ = pool.Exec(ctx, "DELETE FROM auth.users WHERE id = $1", user1.ID)
}

func TestUsersPG_GetBySubject(t *testing.T) {
	pool := openAuthPool(t)
	ctx := context.Background()
	u := &UsersPG{Pool: pool}

	subject := "test|" + uuid.New().String()
	inserted, _ := u.UpsertBySubject(ctx, subject, "get@example.com", "Get User")

	found, err := u.GetBySubject(ctx, subject)
	if err != nil {
		t.Fatalf("GetBySubject: %v", err)
	}
	if found.ID != inserted.ID {
		t.Errorf("ID = %s, want %s", found.ID, inserted.ID)
	}
	if found.Email != "get@example.com" {
		t.Errorf("Email = %q", found.Email)
	}

	_, _ = pool.Exec(ctx, "DELETE FROM auth.users WHERE id = $1", inserted.ID)
}

func TestUsersPG_GetBySubject_NotFound(t *testing.T) {
	pool := openAuthPool(t)
	ctx := context.Background()
	u := &UsersPG{Pool: pool}

	_, err := u.GetBySubject(ctx, "nonexistent|"+uuid.New().String())
	if err == nil {
		t.Fatal("expected error for nonexistent subject")
	}
}
