package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testSecret = "shh-it's-a-secret"

func setup(t *testing.T) (*httptest.Server, *InMemoryRecorder) {
	t.Helper()
	rec := &InMemoryRecorder{}
	h := &Handler{
		Secrets:  func(_, _ string) (string, bool) { return testSecret, true },
		Recorder: rec,
	}
	srv := httptest.NewServer(h.Routes())
	t.Cleanup(srv.Close)
	return srv, rec
}

func sign(body []byte) string {
	m := hmac.New(sha256.New, []byte(testSecret))
	m.Write(body)
	return "sha256=" + hex.EncodeToString(m.Sum(nil))
}

func TestGitHub_AcceptsValidSignature(t *testing.T) {
	srv, rec := setup(t)
	body := []byte(`{"ref":"refs/heads/main","repository":{"full_name":"o/r"}}`)
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/webhooks/github", strings.NewReader(string(body)))
	req.Header.Set("X-Hub-Signature-256", sign(body))
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-GitHub-Delivery", "abc-123")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("want 204 got %d", resp.StatusCode)
	}
	if len(rec.Snapshot()) != 1 {
		t.Fatal("event not recorded")
	}
}

func TestGitHub_RejectsTamper(t *testing.T) {
	srv, _ := setup(t)
	original := []byte(`{"ref":"refs/heads/main"}`)
	tampered := []byte(`{"ref":"refs/heads/evil"}`)
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/webhooks/github", strings.NewReader(string(tampered)))
	req.Header.Set("X-Hub-Signature-256", sign(original))
	req.Header.Set("X-GitHub-Event", "push")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("want 401 got %d", resp.StatusCode)
	}
}

func TestGitHub_RejectsMissingSig(t *testing.T) {
	srv, _ := setup(t)
	resp, err := http.Post(srv.URL+"/v1/webhooks/github", "application/json", strings.NewReader(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("want 401 got %d", resp.StatusCode)
	}
}

func TestGitLab_ChecksToken(t *testing.T) {
	srv, rec := setup(t)
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/webhooks/gitlab", strings.NewReader(`{}`))
	req.Header.Set("X-Gitlab-Token", testSecret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("want 204 got %d", resp.StatusCode)
	}
	if len(rec.Snapshot()) != 1 {
		t.Fatal("not recorded")
	}

	req2, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/webhooks/gitlab", strings.NewReader(`{}`))
	req2.Header.Set("X-Gitlab-Token", "wrong")
	resp2, _ := http.DefaultClient.Do(req2)
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusUnauthorized {
		t.Fatalf("want 401 got %d", resp2.StatusCode)
	}
}

func TestGitea_AcceptsValidSignature(t *testing.T) {
	srv, rec := setup(t)
	body := []byte(`{"ref":"refs/heads/main","repository":{"full_name":"o/r"}}`)
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/webhooks/gitea", strings.NewReader(string(body)))
	req.Header.Set("X-Gitea-Signature", sign(body))
	req.Header.Set("X-Gitea-Event", "push")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("want 204 got %d", resp.StatusCode)
	}
	snapshot := rec.Snapshot()
	if len(snapshot) != 1 || snapshot[0].Provider != "gitea" {
		t.Fatalf("recorded wrong event: %+v", snapshot)
	}
}
