// Package adapter defines the provider-agnostic interface for upstream Git
// forges. Implementations live under internal/providers/<name>/.
package adapter

import (
	"context"
	"time"
)

type Provider string

const (
	GitHub Provider = "github"
	GitLab Provider = "gitlab"
	Gitea  Provider = "gitea"
)

type Source struct {
	Provider Provider
	BaseURL  string
	Token    string
	Owner    string
	Repo     string
}

type Destination = Source

type Branch struct {
	Repo Source
	Name string
}

type RefValue struct {
	Name string
	SHA  string
}

type RefUpdate struct {
	Name   string
	OldSHA string
	NewSHA string
}

type PullRequest struct {
	Number int
	URL    string
}

type RepoInfo struct {
	Default   string
	Private   bool
	UpdatedAt time.Time
}

type Webhook struct {
	ID     string
	URL    string
	Events []string
}

// Adapter is the contract every provider implementation satisfies. M4 ships
// GitHub, GitLab, and Gitea (also powering Codeberg/Forgejo); the remaining
// seven providers land in M5.
type Adapter interface {
	Push(ctx context.Context, dst Destination, refs []RefUpdate) error
	Fetch(ctx context.Context, src Source, refs []string) ([]RefValue, error)
	CreatePR(ctx context.Context, src, dst Branch, title, body string) (*PullRequest, error)
	ListRefs(ctx context.Context, src Source) ([]RefValue, error)
	GetRepo(ctx context.Context, src Source) (*RepoInfo, error)
	ListWebhooks(ctx context.Context, src Source) ([]Webhook, error)
	RegisterWebhook(ctx context.Context, src Source, url, secret string, events []string) (*Webhook, error)
}
