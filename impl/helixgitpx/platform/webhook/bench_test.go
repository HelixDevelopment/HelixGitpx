package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

var benchBody = []byte(`{"ref":"refs/heads/main","repository":{"full_name":"o/r"},"commits":[` +
	`{"id":"abc","message":"first"},{"id":"def","message":"second"}` +
	`]}`)
var benchSecret = []byte("a fairly long and typical webhook secret")

func BenchmarkVerifyHMAC_Valid(b *testing.B) {
	mac := hmac.New(sha256.New, benchSecret)
	mac.Write(benchBody)
	sig := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !VerifyHMAC(benchSecret, benchBody, sig) {
			b.Fatal("expected ok")
		}
	}
}

func BenchmarkVerifyHMAC_Invalid(b *testing.B) {
	bad := "sha256=" + hex.EncodeToString(make([]byte, 32))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = VerifyHMAC(benchSecret, benchBody, bad)
	}
}
