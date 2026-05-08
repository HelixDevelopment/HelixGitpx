package repo

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestMFAPG_InsertTOTPAndGetTOTP(t *testing.T) {
	pool := openAuthPool(t)
	ctx := context.Background()
	users := &UsersPG{Pool: pool}
	mfa := &MFAPG{Pool: pool}

	subject := "test|" + uuid.New().String()
	user, _ := users.UpsertBySubject(ctx, subject, "mfa@example.com", "MFA User")
	defer pool.Exec(ctx, "DELETE FROM auth.users WHERE id = $1", user.ID)

	secret := "JBSWY3DPEHPK3PXP"
	id, err := mfa.InsertTOTP(ctx, user.ID.String(), secret)
	if err != nil {
		t.Fatalf("InsertTOTP: %v", err)
	}
	if id == uuid.Nil {
		t.Error("InsertTOTP returned nil UUID")
	}
	defer pool.Exec(ctx, "DELETE FROM auth.mfa_factors WHERE id = $1", id)

	f, err := mfa.GetTOTP(ctx, user.ID.String())
	if err != nil {
		t.Fatalf("GetTOTP: %v", err)
	}
	if f.ID != id {
		t.Errorf("ID = %s, want %s", f.ID, id)
	}
	if f.UserID != user.ID.String() {
		t.Errorf("UserID = %q, want %q", f.UserID, user.ID.String())
	}
	if f.Kind != "totp" {
		t.Errorf("Kind = %q, want totp", f.Kind)
	}
	if string(f.Secret) != secret {
		t.Errorf("Secret = %q, want %q", string(f.Secret), secret)
	}
}

func TestMFAPG_GetTOTP_NotFound(t *testing.T) {
	pool := openAuthPool(t)
	ctx := context.Background()
	mfa := &MFAPG{Pool: pool}

	_, err := mfa.GetTOTP(ctx, uuid.New().String())
	if err == nil {
		t.Error("expected error for user with no TOTP")
	}
}
