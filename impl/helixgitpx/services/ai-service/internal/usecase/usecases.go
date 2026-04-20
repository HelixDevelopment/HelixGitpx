// Package usecase holds ai-service's business-logic surfaces.
// Each function is a thin wrapper around a LiteLLM router call; the real HTTP
// client is injected in app.Run.
package usecase

import "context"

type Client interface {
	Prompt(ctx context.Context, model, prompt string) (string, error)
}

type UseCases struct {
	LLM Client
}

func (u *UseCases) Summarize(ctx context.Context, content string) (string, error) {
	return u.LLM.Prompt(ctx, "helixgitpx/summarize", "Summarize: "+content)
}

func (u *UseCases) ProposeConflict(ctx context.Context, diff string) (string, error) {
	return u.LLM.Prompt(ctx, "helixgitpx/conflict", "Resolve conflict: "+diff)
}

func (u *UseCases) SuggestLabel(ctx context.Context, title, body string) (string, error) {
	return u.LLM.Prompt(ctx, "helixgitpx/label", "Labels for: "+title+"\n"+body)
}

func (u *UseCases) ChatOps(ctx context.Context, prompt string) (string, error) {
	return u.LLM.Prompt(ctx, "helixgitpx/chatops", prompt)
}
