package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/helixgitpx/platform/auth"
)

func TestSignAndValidate_RoundTrip(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	signer := auth.NewSigner(priv, "kid-1", "helixgitpx")

	tok, err := signer.Issue(auth.Claims{
		Subject: "user-abc",
		Orgs:    []string{"acme"},
		TTL:     15 * time.Minute,
	})
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	validator := auth.NewValidatorFromKey(&priv.PublicKey, "kid-1", "helixgitpx")
	claims, err := validator.Validate(tok)
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if claims.Subject != "user-abc" {
		t.Errorf("subject = %q", claims.Subject)
	}
	if len(claims.Orgs) != 1 || claims.Orgs[0] != "acme" {
		t.Errorf("orgs = %v", claims.Orgs)
	}
}

func TestValidate_ExpiredToken(t *testing.T) {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	signer := auth.NewSigner(priv, "kid-1", "helixgitpx")
	tok, _ := signer.Issue(auth.Claims{Subject: "u", TTL: -1 * time.Second})

	v := auth.NewValidatorFromKey(&priv.PublicKey, "kid-1", "helixgitpx")
	if _, err := v.Validate(tok); err == nil {
		t.Fatal("expected expired error")
	}
}
