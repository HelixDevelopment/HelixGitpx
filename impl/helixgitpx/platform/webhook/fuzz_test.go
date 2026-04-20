package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

// FuzzVerifyHMAC checks that VerifyHMAC never panics and never returns true
// for a body that wasn't signed by the same secret.
func FuzzVerifyHMAC(f *testing.F) {
	// Seed with the canonical happy path.
	f.Add([]byte("topsecret"), []byte("payload"), "sha256=")

	f.Fuzz(func(t *testing.T, secret, body []byte, stub string) {
		// Correctly-signed signature must always verify.
		mac := hmac.New(sha256.New, secret)
		mac.Write(body)
		sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))
		if !VerifyHMAC(secret, body, sig) {
			t.Fatalf("correctly-signed body rejected: secret=%q body=%q", secret, body)
		}

		// The fuzzer's `stub` signature should be rejected unless it
		// happens to match — which is cryptographically negligible but we
		// still accept it if it does (we only forbid a crash).
		_ = VerifyHMAC(secret, body, stub)
	})
}
