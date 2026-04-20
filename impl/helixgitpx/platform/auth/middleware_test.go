package auth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/helixgitpx/platform/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestUnaryInterceptor_InjectsClaims(t *testing.T) {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	signer := auth.NewSigner(priv, "kid-1", "helixgitpx")
	validator := auth.NewValidatorFromKey(&priv.PublicKey, "kid-1", "helixgitpx")

	tok, _ := signer.Issue(auth.Claims{Subject: "user-xyz", TTL: 1 * time.Minute})

	handler := func(ctx context.Context, req any) (any, error) {
		if uid, _ := auth.UserIDFromContext(ctx); uid != "user-xyz" {
			t.Errorf("user_id from context = %q", uid)
		}
		return "ok", nil
	}
	intc := auth.UnaryInterceptor(validator)
	ctx := metadata.NewIncomingContext(context.Background(),
		metadata.Pairs("authorization", "Bearer "+tok))
	if _, err := intc(ctx, nil, &grpc.UnaryServerInfo{}, handler); err != nil {
		t.Fatalf("intc: %v", err)
	}
}

func TestUnaryInterceptor_MissingToken(t *testing.T) {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	validator := auth.NewValidatorFromKey(&priv.PublicKey, "kid-1", "helixgitpx")
	intc := auth.UnaryInterceptor(validator)
	handler := func(ctx context.Context, req any) (any, error) { return nil, nil }
	_, err := intc(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)
	if err == nil {
		t.Fatal("expected unauthenticated error")
	}
}
