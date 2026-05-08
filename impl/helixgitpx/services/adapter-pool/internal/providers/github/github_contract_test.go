package github_test

import (
	"context"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/adapter-pool/internal/adapter"
	provider "github.com/helixgitpx/helixgitpx/services/adapter-pool/internal/providers/github"
)

// TestGitHub_AdapterImplementsInterface verifies the concrete Adapter satisfies
// the adapter.Adapter interface at compile time. This is a real type-safety
// guarantee — if the interface changes, this test won't compile.
var _ adapter.Adapter = (*provider.Adapter)(nil)

// TestGitHub_AdapterMethodsDoNotPanic verifies every adapter method can be
// called without panicking. The GitHub adapter is currently a stub (M4 ships
// the shape; real go-vcr cassette replay arrives with google/go-github SDK
// integration). These tests confirm the stub is callable and returns
// well-formed values.
//
// SKIP-OK: #HGX-M4 — Real contract tests with go-vcr cassettes will replace
// these when the adapter-pool service is wired end-to-end.
func TestGitHub_GetRepo_StubReturnsWellFormedValue(t *testing.T) {
	a := &provider.Adapter{}
	info, err := a.GetRepo(context.Background(), adapter.Source{
		Provider: adapter.GitHub, Owner: "octocat", Repo: "Hello-World",
	})
	if err != nil {
		t.Fatalf("GetRepo: %v", err)
	}
	if info == nil {
		t.Fatal("GetRepo returned nil RepoInfo")
	}
	if info.Default == "" {
		t.Error("default branch must not be empty")
	}
}

func TestGitHub_CreatePR_StubReturnsWellFormedValue(t *testing.T) {
	a := &provider.Adapter{}
	pr, err := a.CreatePR(context.Background(),
		adapter.Branch{Repo: adapter.Source{Provider: adapter.GitHub, Owner: "octocat", Repo: "Hello-World"}, Name: "main"},
		adapter.Branch{Repo: adapter.Source{Provider: adapter.GitHub, Owner: "octocat", Repo: "Hello-World"}, Name: "feature"},
		"Test PR", "body text",
	)
	if err != nil {
		t.Fatalf("CreatePR: %v", err)
	}
	if pr == nil {
		t.Fatal("CreatePR returned nil PullRequest")
	}
}

func TestGitHub_ListWebhooks_StubReturnsNoError(t *testing.T) {
	a := &provider.Adapter{}
	_, err := a.ListWebhooks(context.Background(), adapter.Source{
		Provider: adapter.GitHub, Owner: "octocat", Repo: "Hello-World",
	})
	if err != nil {
		t.Fatalf("ListWebhooks: %v", err)
	}
}

func TestGitHub_RegisterWebhook_StubReturnsWellFormedValue(t *testing.T) {
	a := &provider.Adapter{}
	wh, err := a.RegisterWebhook(context.Background(), adapter.Source{
		Provider: adapter.GitHub, Owner: "octocat", Repo: "Hello-World",
	}, "https://example.com/hook", "secret123", []string{"push"})
	if err != nil {
		t.Fatalf("RegisterWebhook: %v", err)
	}
	if wh == nil {
		t.Fatal("RegisterWebhook returned nil Webhook")
	}
	if wh.URL != "https://example.com/hook" {
		t.Errorf("Webhook URL = %q, want %q", wh.URL, "https://example.com/hook")
	}
}

func TestGitHub_Push_StubReturnsNoError(t *testing.T) {
	a := &provider.Adapter{}
	err := a.Push(context.Background(), adapter.Source{Provider: adapter.GitHub}, []adapter.RefUpdate{{Name: "refs/heads/main", OldSHA: "abc", NewSHA: "def"}})
	if err != nil {
		t.Fatalf("Push: %v", err)
	}
}

func TestGitHub_Fetch_StubReturnsNoError(t *testing.T) {
	a := &provider.Adapter{}
	_, err := a.Fetch(context.Background(), adapter.Source{Provider: adapter.GitHub, Owner: "octocat", Repo: "Hello-World"}, []string{"refs/heads/main"})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
}

func TestGitHub_ListRefs_StubReturnsNoError(t *testing.T) {
	a := &provider.Adapter{}
	_, err := a.ListRefs(context.Background(), adapter.Source{
		Provider: adapter.GitHub, Owner: "octocat", Repo: "Hello-World",
	})
	if err != nil {
		t.Fatalf("ListRefs: %v", err)
	}
}
