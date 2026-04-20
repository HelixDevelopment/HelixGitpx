package canonical_test

import (
	"testing"

	"github.com/helixgitpx/helixgitpx/services/webhook-gateway/internal/canonical"
)

func TestCanonicalizeGitHub(t *testing.T) {
	body := []byte(`{"ref":"refs/heads/main","repository":{"full_name":"acme/example"}}`)
	ev := canonical.CanonicalizeGitHub("del-123", "push", body)
	if ev.Provider != "github" || ev.DeliveryID != "del-123" || ev.EventType != "push" {
		t.Errorf("canonicalize: %+v", ev)
	}
	if ev.Repo != "acme/example" || ev.Ref != "refs/heads/main" {
		t.Errorf("repo/ref extraction failed: %+v", ev)
	}
}

func TestCanonicalizeGitLab(t *testing.T) {
	body := []byte(`{"ref":"refs/heads/main","project":{"path_with_namespace":"acme/example"}}`)
	ev := canonical.CanonicalizeGitLab("abc", "Push Hook", body)
	if ev.Repo != "acme/example" {
		t.Errorf("gitlab repo: %q", ev.Repo)
	}
}

func TestCanonicalizeGitea(t *testing.T) {
	body := []byte(`{"ref":"refs/heads/main","repository":{"full_name":"acme/example"}}`)
	ev := canonical.CanonicalizeGitea("xyz", "push", body)
	if ev.Provider != "gitea" || ev.Repo != "acme/example" {
		t.Errorf("gitea canonicalize: %+v", ev)
	}
}
