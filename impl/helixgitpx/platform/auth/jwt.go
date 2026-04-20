package auth

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is the HelixGitpx JWT payload (subset of RFC 7519 + custom).
type Claims struct {
	Subject string        // user id
	Orgs    []string      // org slugs the user is in
	TTL     time.Duration // token lifetime; negative for already-expired (test)
}

// Signer issues RS256 JWTs for a single key.
type Signer struct {
	priv     *rsa.PrivateKey
	kid      string
	issuer   string
	audience string
}

// NewSigner constructs a Signer. audience defaults to issuer.
func NewSigner(priv *rsa.PrivateKey, kid, issuer string) *Signer {
	return &Signer{priv: priv, kid: kid, issuer: issuer, audience: issuer}
}

// PublicKey exposes the signer's RSA public key for JWKS publication.
func (s *Signer) PublicKey() *rsa.PublicKey { return &s.priv.PublicKey }

// KID returns the key id used in the JWT header.
func (s *Signer) KID() string { return s.kid }

// Issue mints a JWT with the given claims.
func (s *Signer) Issue(c Claims) (string, error) {
	now := time.Now()
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":  c.Subject,
		"orgs": c.Orgs,
		"iss":  s.issuer,
		"aud":  s.audience,
		"iat":  now.Unix(),
		"exp":  now.Add(c.TTL).Unix(),
	})
	tok.Header["kid"] = s.kid
	signed, err := tok.SignedString(s.priv)
	if err != nil {
		return "", fmt.Errorf("auth: sign: %w", err)
	}
	return signed, nil
}

// Validator verifies RS256 JWTs against a known public key.
type Validator struct {
	pub      *rsa.PublicKey
	kid      string
	issuer   string
	audience string
}

// NewValidatorFromKey constructs a Validator from a static public key.
func NewValidatorFromKey(pub *rsa.PublicKey, kid, issuer string) *Validator {
	return &Validator{pub: pub, kid: kid, issuer: issuer, audience: issuer}
}

// Validate parses and verifies token, returning Claims or an error.
func (v *Validator) Validate(token string) (Claims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("auth: unexpected signing method: %v", t.Header["alg"])
		}
		if kid, _ := t.Header["kid"].(string); kid != v.kid {
			return nil, fmt.Errorf("auth: kid mismatch")
		}
		return v.pub, nil
	})
	if err != nil {
		return Claims{}, fmt.Errorf("auth: validate: %w", err)
	}
	mc, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || !parsed.Valid {
		return Claims{}, fmt.Errorf("auth: invalid token")
	}
	if iss, _ := mc["iss"].(string); iss != v.issuer {
		return Claims{}, fmt.Errorf("auth: issuer mismatch")
	}
	out := Claims{Subject: asString(mc["sub"])}
	if orgs, ok := mc["orgs"].([]any); ok {
		for _, o := range orgs {
			if s, ok := o.(string); ok {
				out.Orgs = append(out.Orgs, s)
			}
		}
	}
	return out, nil
}

func asString(v any) string {
	s, _ := v.(string)
	return s
}
