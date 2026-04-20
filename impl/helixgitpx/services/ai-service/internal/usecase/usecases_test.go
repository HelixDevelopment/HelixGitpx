package usecase_test

import (
	"context"
	"strings"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/ai-service/internal/usecase"
)

type fakeLLM struct{ last string }

func (f *fakeLLM) Prompt(_ context.Context, _, prompt string) (string, error) {
	f.last = prompt
	return "ok:" + prompt, nil
}

func TestUseCases_Summarize(t *testing.T) {
	llm := &fakeLLM{}
	uc := &usecase.UseCases{LLM: llm}
	got, err := uc.Summarize(context.Background(), "x")
	if err != nil {
		t.Fatalf("Summarize: %v", err)
	}
	if !strings.HasPrefix(got, "ok:") {
		t.Errorf("got %q", got)
	}
}

func TestUseCases_Labels(t *testing.T) {
	llm := &fakeLLM{}
	uc := &usecase.UseCases{LLM: llm}
	_, _ = uc.SuggestLabel(context.Background(), "title", "body")
	if !strings.Contains(llm.last, "Labels for:") {
		t.Errorf("label prompt missing expected header: %q", llm.last)
	}
}
