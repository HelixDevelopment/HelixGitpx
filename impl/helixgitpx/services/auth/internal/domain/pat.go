// Package domain holds the auth service business logic.
// pat.go: HelixGitpx Personal Access Tokens — format hpxat_ + base62(24B),
// stored as SHA-256 hash of the full token string.
package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"math/big"
	"strings"
)

const (
	patPrefix = "hpxat_"
	patBytes  = 24
	patWidth  = 32
)

var base62alphabet = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

// IssuePAT returns a token (returned to user once) and its hashed form (for storage).
func IssuePAT() (token string, hashed []byte, err error) {
	raw := make([]byte, patBytes)
	if _, err = rand.Read(raw); err != nil {
		return "", nil, fmt.Errorf("pat: read random: %w", err)
	}
	encoded := base62encode(raw)
	token = patPrefix + encoded
	h := sha256.Sum256([]byte(token))
	return token, h[:], nil
}

// VerifyPAT returns true iff the presented token hashes to the stored digest.
func VerifyPAT(token string, hashed []byte) bool {
	if !strings.HasPrefix(token, patPrefix) {
		return false
	}
	h := sha256.Sum256([]byte(token))
	return subtle.ConstantTimeCompare(h[:], hashed) == 1
}

func base62encode(buf []byte) string {
	n := new(big.Int).SetBytes(buf)
	base := big.NewInt(62)
	var out []byte
	for n.Sign() > 0 {
		mod := new(big.Int)
		n.DivMod(n, base, mod)
		out = append([]byte{base62alphabet[mod.Int64()]}, out...)
	}
	for len(out) < patWidth {
		out = append([]byte{base62alphabet[0]}, out...)
	}
	return string(out)
}
