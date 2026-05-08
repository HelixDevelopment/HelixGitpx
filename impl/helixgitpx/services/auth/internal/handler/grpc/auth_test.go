package grpc

import (
	"context"
	"testing"

	pb "github.com/helixgitpx/helixgitpx/gen/go/helixgitpx/auth/v1"
	hauth "github.com/helixgitpx/platform/auth"
)

func ctxWithUser(uid string) context.Context {
	return hauth.ContextWithUserID(context.Background(), uid)
}

func TestServer_IssuePAT_NoUserInContext(t *testing.T) {
	srv := &Server{}
	_, err := srv.IssuePAT(context.Background(), &pb.IssuePATRequest{Name: "test"})
	if err == nil {
		t.Fatal("expected error for missing user in context")
	}
}

func TestServer_ListPATs_NilDeps_Panics(t *testing.T) {
	srv := &Server{}
	ctx := ctxWithUser("user-1")
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic with nil PATs repo")
		}
	}()
	srv.ListPATs(ctx, nil)
}

func TestServer_RevokePAT_NilDeps_Panics(t *testing.T) {
	srv := &Server{}
	ctx := ctxWithUser("user-1")
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic with nil PATs repo")
		}
	}()
	srv.RevokePAT(ctx, &pb.RevokePATRequest{Id: "pat-1"})
}

func TestServer_ListSessions_NilDeps_Panics(t *testing.T) {
	srv := &Server{}
	ctx := ctxWithUser("user-1")
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic with nil Sessions repo")
		}
	}()
	srv.ListSessions(ctx, nil)
}

func TestServer_RevokeSession_InvalidID(t *testing.T) {
	srv := &Server{}
	ctx := ctxWithUser("user-1")
	_, err := srv.RevokeSession(ctx, &pb.RevokeSessionRequest{Id: "not-a-uuid"})
	if err == nil {
		t.Fatal("expected error for invalid session ID")
	}
}

func TestServer_EnrollTOTP_NilDeps_Panics(t *testing.T) {
	srv := &Server{}
	ctx := ctxWithUser("user-1")
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic with nil Users repo")
		}
	}()
	srv.EnrollTOTP(ctx, nil)
}

func TestServer_VerifyMFA_NoCode(t *testing.T) {
	srv := &Server{}
	ctx := ctxWithUser("user-1")
	resp, err := srv.VerifyMFA(ctx, &pb.VerifyMFARequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Verified {
		t.Error("should not be verified with no code")
	}
}

func TestServer_VerifyMFA_TOTPNotEnrolled(t *testing.T) {
	srv := &Server{}
	ctx := ctxWithUser("user-1")
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic with nil MFA repo")
		}
	}()
	srv.VerifyMFA(ctx, &pb.VerifyMFARequest{TotpCode: "123456"})
}

func TestServer_ExchangeOIDC_Unimplemented(t *testing.T) {
	srv := &Server{}
	_, err := srv.ExchangeOIDC(context.Background(), &pb.ExchangeOIDCRequest{})
	if err == nil {
		t.Fatal("expected Unimplemented error")
	}
}

func TestServer_RefreshToken_Unimplemented(t *testing.T) {
	srv := &Server{}
	_, err := srv.RefreshToken(context.Background(), &pb.RefreshTokenRequest{})
	if err == nil {
		t.Fatal("expected Unimplemented error")
	}
}

func TestServer_EnrollFIDO2_Unimplemented(t *testing.T) {
	srv := &Server{}
	_, err := srv.EnrollFIDO2(context.Background(), &pb.EnrollFIDO2Request{})
	if err == nil {
		t.Fatal("expected Unimplemented error")
	}
}

func TestServer_WhoAmI_NilDeps_Panics(t *testing.T) {
	srv := &Server{}
	ctx := ctxWithUser("user-1")
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic with nil Users repo")
		}
	}()
	srv.WhoAmI(ctx, nil)
}

func TestServer_WhoAmI_NoUserInContext_Panics(t *testing.T) {
	srv := &Server{}
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic with no user in context and nil Users repo")
		}
	}()
	srv.WhoAmI(context.Background(), nil)
}
