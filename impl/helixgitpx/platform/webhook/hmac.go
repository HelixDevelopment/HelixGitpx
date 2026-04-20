package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// VerifyHMAC returns true iff signature (in "sha256=HEX" form) is a correct
// HMAC-SHA256 of body under secret. Constant-time compare.
func VerifyHMAC(secret, body []byte, signature string) bool {
	signature = strings.TrimPrefix(signature, "sha256=")
	want, err := hex.DecodeString(signature)
	if err != nil || len(want) == 0 {
		return false
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	got := mac.Sum(nil)
	return hmac.Equal(got, want)
}
