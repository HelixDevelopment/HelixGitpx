package repo

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func createTestUser(t *testing.T, ctx context.Context, pool interface{ Exec(ctx context.Context, sql string, args ...any) (interface{}, error) }) string {
	t.Helper()
	return ""
}

func TestPATsPG_InsertAndList(t *testing.T) {
	pool := openAuthPool(t)
	ctx := context.Background()
	users := &UsersPG{Pool: pool}
	pats := &PATsPG{Pool: pool}

	subject := "test|" + uuid.New().String()
	user, _ := users.UpsertBySubject(ctx, subject, "pat@example.com", "PAT User")
	defer pool.Exec(ctx, "DELETE FROM auth.users WHERE id = $1", user.ID)

	scopes := []string{"repo:read", "repo:write"}
	secret := []byte("hashed-secret-abc")
	pat, err := pats.Insert(ctx, user.ID.String(), "my-token", secret, scopes, nil)
	if err != nil {
		t.Fatalf("Insert: %v", err)
	}
	if pat.ID == uuid.Nil {
		t.Error("pat ID should not be nil UUID")
	}
	if pat.Name != "my-token" {
		t.Errorf("Name = %q, want my-token", pat.Name)
	}
	if pat.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}

	list, err := pats.List(ctx, user.ID.String())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	found := false
	for _, p := range list {
		if p.ID == pat.ID {
			found = true
			if p.Name != "my-token" {
				t.Errorf("Name = %q, want my-token", p.Name)
			}
			if len(p.Scopes) != 2 || p.Scopes[0] != "repo:read" || p.Scopes[1] != "repo:write" {
				t.Errorf("Scopes = %v, want [repo:read repo:write]", p.Scopes)
			}
		}
	}
	if !found {
		t.Error("inserted PAT not found in List results")
	}

	_, _ = pool.Exec(ctx, "DELETE FROM auth.pats WHERE id = $1", pat.ID)
}

func TestPATsPG_RevokeExcludesFromList(t *testing.T) {
	pool := openAuthPool(t)
	ctx := context.Background()
	users := &UsersPG{Pool: pool}
	pats := &PATsPG{Pool: pool}

	subject := "test|" + uuid.New().String()
	user, _ := users.UpsertBySubject(ctx, subject, "revoke-pat@example.com", "Revoke PAT User")
	defer pool.Exec(ctx, "DELETE FROM auth.users WHERE id = $1", user.ID)

	pat, _ := pats.Insert(ctx, user.ID.String(), "doomed", []byte("s"), []string{"read"}, nil)
	defer pool.Exec(ctx, "DELETE FROM auth.pats WHERE id = $1", pat.ID)

	if err := pats.Revoke(ctx, pat.ID.String(), user.ID.String()); err != nil {
		t.Fatalf("Revoke: %v", err)
	}

	list, err := pats.List(ctx, user.ID.String())
	if err != nil {
		t.Fatalf("List after revoke: %v", err)
	}
	for _, p := range list {
		if p.ID == pat.ID {
			t.Error("revoked PAT should not appear in List results")
		}
	}
}

func TestPATsPG_InsertWithExpiry(t *testing.T) {
	pool := openAuthPool(t)
	ctx := context.Background()
	users := &UsersPG{Pool: pool}
	pats := &PATsPG{Pool: pool}

	subject := "test|" + uuid.New().String()
	user, _ := users.UpsertBySubject(ctx, subject, "expiry-pat@example.com", "Expiry User")
	defer pool.Exec(ctx, "DELETE FROM auth.users WHERE id = $1", user.ID)

	expires := time.Now().Add(24 * time.Hour)
	pat, err := pats.Insert(ctx, user.ID.String(), "expiring", []byte("s"), nil, &expires)
	if err != nil {
		t.Fatalf("Insert with expiry: %v", err)
	}
	if pat.ExpiresAt == nil {
		t.Fatal("ExpiresAt should be set when provided")
	}

	list, _ := pats.List(ctx, user.ID.String())
	for _, p := range list {
		if p.ID == pat.ID {
			if p.ExpiresAt == nil {
				t.Error("List should return ExpiresAt")
			}
		}
	}

	_, _ = pool.Exec(ctx, "DELETE FROM auth.pats WHERE id = $1", pat.ID)
}
