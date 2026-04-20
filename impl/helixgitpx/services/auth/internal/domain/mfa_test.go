package domain_test

import (
	"testing"
	"time"

	"github.com/helixgitpx/helixgitpx/services/auth/internal/domain"
	"github.com/pquerna/otp/totp"
)

func TestEnrollTOTP_ThenVerify(t *testing.T) {
	otpauth, secret, err := domain.EnrollTOTP("user@helixgitpx.local")
	if err != nil {
		t.Fatalf("EnrollTOTP: %v", err)
	}
	if len(otpauth) == 0 || len(secret) == 0 {
		t.Fatalf("empty enrollment output")
	}
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		t.Fatalf("GenerateCode: %v", err)
	}
	if !domain.VerifyTOTP(secret, code) {
		t.Errorf("VerifyTOTP rejected a freshly generated code")
	}
	if domain.VerifyTOTP(secret, "000000") {
		t.Errorf("VerifyTOTP accepted clearly-wrong code")
	}
}

func TestFIDO2RelyingPartyID(t *testing.T) {
	if got := domain.FIDO2RelyingPartyID("helixgitpx.local"); got != "helixgitpx.local" {
		t.Errorf("RPID = %q", got)
	}
}
