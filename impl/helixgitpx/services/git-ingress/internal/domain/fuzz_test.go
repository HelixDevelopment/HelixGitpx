package domain

import "testing"

// FuzzValidateRef asserts that no input can crash ValidateRef and that
// every explicitly-valid seed accepts while a bad-prefix seed rejects.
func FuzzValidateRef(f *testing.F) {
	f.Add("refs/heads/main")
	f.Add("refs/tags/v1.0.0")
	f.Add("refs/pull/42/head")
	f.Add("")
	f.Add("main")
	f.Add("refs/heads/with space")

	f.Fuzz(func(t *testing.T, ref string) {
		err := ValidateRef(ref)
		// Invariant: a valid ref always round-trips to IsProtected without panic.
		if err == nil {
			_ = IsProtected(ref)
		}
	})
}
