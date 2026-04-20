// Package gitflic implements adapter.Adapter for the gitflic provider.
// M5 ships a stub satisfying the interface; full API wiring lands when
// the first real upstream of this type is configured (tracked via the
// upstream-service.Create flow with provider=PROVIDER_GITFLIC).
package gitflic

import (
	"context"

	"github.com/helixgitpx/helixgitpx/services/adapter-pool/internal/adapter"
)

type Adapter struct{}

func (a *Adapter) Push(_ context.Context, _ adapter.Destination, _ []adapter.RefUpdate) error { return nil }
func (a *Adapter) Fetch(_ context.Context, _ adapter.Source, _ []string) ([]adapter.RefValue, error) { return nil, nil }
func (a *Adapter) CreatePR(_ context.Context, _, _ adapter.Branch, _, _ string) (*adapter.PullRequest, error) { return &adapter.PullRequest{}, nil }
func (a *Adapter) ListRefs(_ context.Context, _ adapter.Source) ([]adapter.RefValue, error) { return nil, nil }
func (a *Adapter) GetRepo(_ context.Context, _ adapter.Source) (*adapter.RepoInfo, error) { return &adapter.RepoInfo{Default: "main"}, nil }
func (a *Adapter) ListWebhooks(_ context.Context, _ adapter.Source) ([]adapter.Webhook, error) { return nil, nil }
func (a *Adapter) RegisterWebhook(_ context.Context, _ adapter.Source, url, _ string, events []string) (*adapter.Webhook, error) {
	return &adapter.Webhook{URL: url, Events: events}, nil
}
