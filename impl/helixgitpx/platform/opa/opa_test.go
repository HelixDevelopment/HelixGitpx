package opa_test

import (
	"context"
	"testing"

	"github.com/helixgitpx/platform/opa"
)

func TestEval_AllowFromModule(t *testing.T) {
	ev, err := opa.NewEvaluator(context.Background(), opa.Options{
		Module: `package helixgitpx
allow { input.role == "admin" }`,
		Query: "data.helixgitpx.allow",
	})
	if err != nil {
		t.Fatalf("NewEvaluator: %v", err)
	}
	got, err := ev.Eval(context.Background(), map[string]any{"role": "admin"})
	if err != nil {
		t.Fatalf("Eval: %v", err)
	}
	if b, _ := got.(bool); !b {
		t.Errorf("allow = %v, want true", got)
	}
}
