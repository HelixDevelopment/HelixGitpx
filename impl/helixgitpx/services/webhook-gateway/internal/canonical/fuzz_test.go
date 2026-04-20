package canonical

import "testing"

// FuzzCanonicalizeGitHub asserts that no malformed body can crash the
// canonicalizer (all Unmarshal errors must be swallowed) and that the
// resulting Event always carries the original body unchanged.
func FuzzCanonicalizeGitHub(f *testing.F) {
	f.Add("abc-123", "push", []byte(`{"ref":"refs/heads/main","repository":{"full_name":"o/r"}}`))
	f.Add("", "", []byte(``))
	f.Add("d", "x", []byte(`{"ref":null}`))
	f.Add("d", "x", []byte(`not json`))

	f.Fuzz(func(t *testing.T, deliveryID, eventType string, body []byte) {
		evt := CanonicalizeGitHub(deliveryID, eventType, body)
		if evt.Provider != "github" {
			t.Fatalf("provider drift: %q", evt.Provider)
		}
		if string(evt.BodyRaw) != string(body) {
			t.Fatalf("body mutated: got %q want %q", evt.BodyRaw, body)
		}
		if evt.DeliveryID != deliveryID {
			t.Fatalf("delivery id drift")
		}
	})
}

func FuzzCanonicalizeGitLab(f *testing.F) {
	f.Add("abc", "push", []byte(`{"ref":"refs/heads/main","project":{"path_with_namespace":"g/p"}}`))
	f.Fuzz(func(t *testing.T, id, ev string, body []byte) {
		e := CanonicalizeGitLab(id, ev, body)
		if e.Provider != "gitlab" {
			t.Fatal("provider drift")
		}
	})
}

func FuzzCanonicalizeGitea(f *testing.F) {
	f.Add("abc", "push", []byte(`{"ref":"refs/heads/main","repository":{"full_name":"o/r"}}`))
	f.Fuzz(func(t *testing.T, id, ev string, body []byte) {
		e := CanonicalizeGitea(id, ev, body)
		if e.Provider != "gitea" {
			t.Fatal("provider drift")
		}
	})
}
