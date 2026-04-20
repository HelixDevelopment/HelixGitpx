// Package github implements adapter.Adapter for GitHub REST/GraphQL.
// Real API calls are stubbed in M4 — full implementations (Push, Fetch, CreatePR)
// land when adapter-pool integrates with the test cassettes in M5 hardening.
package github

import (
	"context"

	"github.com/helixgitpx/helixgitpx/services/adapter-pool/internal/adapter"
)

// Adapter is the GitHub provider. Real implementations use google/go-github
// via an oauth2 HTTPClient; tests inject a go-vcr transport.
type Adapter struct{}

func (a *Adapter) Push(_ context.Context, _ adapter.Destination, _ []adapter.RefUpdate) error {
	return nil
}
func (a *Adapter) Fetch(_ context.Context, _ adapter.Source, _ []string) ([]adapter.RefValue, error) {
	return nil, nil
}
func (a *Adapter) CreatePR(_ context.Context, _, _ adapter.Branch, _, _ string) (*adapter.PullRequest, error) {
	return &adapter.PullRequest{}, nil
}
func (a *Adapter) ListRefs(_ context.Context, _ adapter.Source) ([]adapter.RefValue, error) {
	return nil, nil
}
func (a *Adapter) GetRepo(_ context.Context, _ adapter.Source) (*adapter.RepoInfo, error) {
	return &adapter.RepoInfo{Default: "main"}, nil
}
func (a *Adapter) ListWebhooks(_ context.Context, _ adapter.Source) ([]adapter.Webhook, error) {
	return nil, nil
}
func (a *Adapter) RegisterWebhook(_ context.Context, _ adapter.Source, url, _ string, events []string) (*adapter.Webhook, error) {
	return &adapter.Webhook{URL: url, Events: events}, nil
}
