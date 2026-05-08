package repo

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSessionsPG_CreateAndActive(t *testing.T) {
	pool := openAuthPool(t)
	ctx := context.Background()
	users := &UsersPG{Pool: pool}
	sessions := &SessionsPG{Pool: pool}

	subject := "test|" + uuid.New().String()
	user, _ := users.UpsertBySubject(ctx, subject, "sess@example.com", "Sess User")
	defer pool.Exec(ctx, "DELETE FROM auth.users WHERE id = $1", user.ID)

	sid := uuid.New()
	expires := time.Now().Add(1 * time.Hour)
	if err := sessions.Create(ctx, sid, user.ID.String(), expires, "test-agent", "127.0.0.1"); err != nil {
		t.Fatalf("Create: %v", err)
	}

	sess, err := sessions.Active(ctx, sid)
	if err != nil {
		t.Fatalf("Active: %v", err)
	}
	if sess == nil {
		t.Fatal("Active returned nil for active session")
	}
	if sess.UserID != user.ID.String() {
		t.Errorf("UserID = %q, want %q", sess.UserID, user.ID.String())
	}
	if sess.UserAgent != "test-agent" {
		t.Errorf("UserAgent = %q", sess.UserAgent)
	}
}

func TestSessionsPG_RevokeThenInactive(t *testing.T) {
	pool := openAuthPool(t)
	ctx := context.Background()
	users := &UsersPG{Pool: pool}
	sessions := &SessionsPG{Pool: pool}

	subject := "test|" + uuid.New().String()
	user, _ := users.UpsertBySubject(ctx, subject, "revoke@example.com", "Revoke User")
	defer pool.Exec(ctx, "DELETE FROM auth.users WHERE id = $1", user.ID)

	sid := uuid.New()
	expires := time.Now().Add(1 * time.Hour)
	_ = sessions.Create(ctx, sid, user.ID.String(), expires, "", "")

	if err := sessions.Revoke(ctx, sid, user.ID.String()); err != nil {
		t.Fatalf("Revoke: %v", err)
	}

	sess, err := sessions.Active(ctx, sid)
	if err != nil {
		t.Fatalf("Active after revoke: %v", err)
	}
	if sess != nil {
		t.Error("revoked session should not be active")
	}
}

func TestSessionsPG_List(t *testing.T) {
	pool := openAuthPool(t)
	ctx := context.Background()
	users := &UsersPG{Pool: pool}
	sessions := &SessionsPG{Pool: pool}

	subject := "test|" + uuid.New().String()
	user, _ := users.UpsertBySubject(ctx, subject, "list@example.com", "List User")
	defer pool.Exec(ctx, "DELETE FROM auth.users WHERE id = $1", user.ID)

	sid := uuid.New()
	_ = sessions.Create(ctx, sid, user.ID.String(), time.Now().Add(1*time.Hour), "list-agent", "")

	rows, err := sessions.List(ctx, user.ID.String())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	found := false
	for _, s := range rows {
		if s.ID == sid {
			found = true
			if s.UserAgent != "list-agent" {
				t.Errorf("UserAgent = %q, want list-agent", s.UserAgent)
			}
		}
	}
	if !found {
		t.Error("created session not found in List results")
	}
}
