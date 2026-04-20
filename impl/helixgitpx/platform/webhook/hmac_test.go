package webhook_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/helixgitpx/platform/webhook"
)

func TestVerifyHMAC_GitHubStyle(t *testing.T) {
	secret := []byte("s3cr3t")
	body := []byte(`{"action":"opened"}`)
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	if !webhook.VerifyHMAC(secret, body, sig) {
		t.Errorf("VerifyHMAC rejected a correct signature")
	}
	if webhook.VerifyHMAC(secret, body, "sha256=00000000") {
		t.Errorf("VerifyHMAC accepted a wrong signature")
	}
	if webhook.VerifyHMAC(secret, body, "") {
		t.Errorf("VerifyHMAC accepted empty signature")
	}
}
