// Package domain (mfa) handles TOTP + FIDO2 enrollment and verification.
// TOTP lives here; FIDO2 requires a live HTTP request context and is wired
// at the handler/http layer.
package domain

import "github.com/pquerna/otp/totp"

// EnrollTOTP generates a TOTP secret for account. Returns the otpauth URL
// (scan as QR) and the raw secret (store in auth.mfa_factors).
func EnrollTOTP(account string) (otpauthURL, secret string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "HelixGitpx",
		AccountName: account,
	})
	if err != nil {
		return "", "", err
	}
	return key.URL(), key.Secret(), nil
}

// VerifyTOTP returns true iff code is valid for secret at the current time.
func VerifyTOTP(secret, code string) bool {
	return totp.Validate(code, secret)
}

// FIDO2RelyingPartyID returns the canonical RPID for the given trust domain.
func FIDO2RelyingPartyID(trustDomain string) string { return trustDomain }
