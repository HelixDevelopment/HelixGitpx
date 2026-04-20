package github_test

import (
	"context"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/adapter-pool/internal/adapter"
	provider "github.com/helixgitpx/helixgitpx/services/adapter-pool/internal/providers/github"
)

// Contract test stub: M4 ships the shape; real go-vcr cassette replay
// arrives when adapter-pool integrates with the actual google/go-github SDK.
func TestGitHub_GetRepo_Contract(t *testing.T) {
	a := &provider.Adapter{}
	info, err := a.GetRepo(context.Background(), adapter.Source{
		Provider: adapter.GitHub, Owner: "octocat", Repo: "Hello-World",
	})
	if err != nil {
		t.Fatalf("GetRepo: %v", err)
	}
	if info.Default == "" {
		t.Errorf("default branch empty")
	}
}
