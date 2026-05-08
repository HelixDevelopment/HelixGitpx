package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/ai-service/internal/usecase"
)

type fakeLLM struct {
	lastModel  string
	lastPrompt string
	err        error
}

func (f *fakeLLM) Prompt(_ context.Context, model, prompt string) (string, error) {
	f.lastModel = model
	f.lastPrompt = prompt
	if f.err != nil {
		return "", f.err
	}
	return "ok:" + prompt, nil
}

func TestUseCases_Summarize_PromptFormat(t *testing.T) {
	llm := &fakeLLM{}
	uc := &usecase.UseCases{LLM: llm}
	got, err := uc.Summarize(context.Background(), "some diff content")
	if err != nil {
		t.Fatalf("Summarize: %v", err)
	}
	if llm.lastModel != "helixgitpx/summarize" {
		t.Errorf("model = %q, want %q", llm.lastModel, "helixgitpx/summarize")
	}
	if !strings.HasPrefix(llm.lastPrompt, "Summarize: ") {
		t.Errorf("prompt = %q, want prefix %q", llm.lastPrompt, "Summarize: ")
	}
	if !strings.Contains(got, "ok:") {
		t.Errorf("response = %q, want ok: prefix", got)
	}
}

func TestUseCases_SuggestLabel_PromptFormat(t *testing.T) {
	llm := &fakeLLM{}
	uc := &usecase.UseCases{LLM: llm}
	_, err := uc.SuggestLabel(context.Background(), "Fix bug", "body text")
	if err != nil {
		t.Fatal(err)
	}
	if llm.lastModel != "helixgitpx/label" {
		t.Errorf("model = %q, want %q", llm.lastModel, "helixgitpx/label")
	}
	if !strings.Contains(llm.lastPrompt, "Labels for: Fix bug") {
		t.Errorf("prompt missing title: %q", llm.lastPrompt)
	}
	if !strings.Contains(llm.lastPrompt, "body text") {
		t.Errorf("prompt missing body: %q", llm.lastPrompt)
	}
}

func TestUseCases_ProposeConflict_PromptFormat(t *testing.T) {
	llm := &fakeLLM{}
	uc := &usecase.UseCases{LLM: llm}
	got, err := uc.ProposeConflict(context.Background(), "<<<< HEAD\nfoo\n====\nbar\n>>>>")
	if err != nil {
		t.Fatal(err)
	}
	if llm.lastModel != "helixgitpx/conflict" {
		t.Errorf("model = %q, want %q", llm.lastModel, "helixgitpx/conflict")
	}
	if !strings.HasPrefix(llm.lastPrompt, "Resolve conflict: ") {
		t.Errorf("prompt = %q, want prefix %q", llm.lastPrompt, "Resolve conflict: ")
	}
	if got == "" {
		t.Error("response is empty")
	}
}

func TestUseCases_ChatOps_PromptFormat(t *testing.T) {
	llm := &fakeLLM{}
	uc := &usecase.UseCases{LLM: llm}
	got, err := uc.ChatOps(context.Background(), "deploy staging")
	if err != nil {
		t.Fatal(err)
	}
	if llm.lastModel != "helixgitpx/chatops" {
		t.Errorf("model = %q, want %q", llm.lastModel, "helixgitpx/chatops")
	}
	if llm.lastPrompt != "deploy staging" {
		t.Errorf("prompt = %q, want %q", llm.lastPrompt, "deploy staging")
	}
	if got == "" {
		t.Error("response is empty")
	}
}

func TestUseCases_Summarize_PropagatesError(t *testing.T) {
	llm := &fakeLLM{err: errors.New("llm: rate limit")}
	uc := &usecase.UseCases{LLM: llm}
	_, err := uc.Summarize(context.Background(), "x")
	if err == nil {
		t.Fatal("expected error from LLM")
	}
	if err.Error() != "llm: rate limit" {
		t.Errorf("err = %q, want %q", err.Error(), "llm: rate limit")
	}
}

func TestUseCases_ProposeConflict_PropagatesError(t *testing.T) {
	llm := &fakeLLM{err: errors.New("llm: timeout")}
	uc := &usecase.UseCases{LLM: llm}
	_, err := uc.ProposeConflict(context.Background(), "diff")
	if err == nil {
		t.Fatal("expected error from LLM")
	}
}

func TestUseCases_SuggestLabel_PropagatesError(t *testing.T) {
	llm := &fakeLLM{err: errors.New("llm: server error")}
	uc := &usecase.UseCases{LLM: llm}
	_, err := uc.SuggestLabel(context.Background(), "t", "b")
	if err == nil {
		t.Fatal("expected error from LLM")
	}
}

func TestUseCases_ChatOps_PropagatesError(t *testing.T) {
	llm := &fakeLLM{err: errors.New("llm: overloaded")}
	uc := &usecase.UseCases{LLM: llm}
	_, err := uc.ChatOps(context.Background(), "deploy")
	if err == nil {
		t.Fatal("expected error from LLM")
	}
}
