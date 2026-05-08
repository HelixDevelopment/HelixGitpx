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

func TestEval_DenyNonAdmin(t *testing.T) {
	ev, err := opa.NewEvaluator(context.Background(), opa.Options{
		Module: `package helixgitpx
allow { input.role == "admin" }`,
		Query: "data.helixgitpx.allow",
	})
	if err != nil {
		t.Fatal(err)
	}
	got, err := ev.Eval(context.Background(), map[string]any{"role": "viewer"})
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := got.(bool); b {
		t.Errorf("allow = %v, want false for non-admin", got)
	}
}

func TestEval_EmptyResultSet(t *testing.T) {
	ev, err := opa.NewEvaluator(context.Background(), opa.Options{
		Module: `package helixgitpx`,
		Query:  "data.helixgitpx.nonexistent",
	})
	if err != nil {
		t.Fatal(err)
	}
	got, err := ev.Eval(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != false {
		t.Errorf("want false for empty result, got %v", got)
	}
}

func TestNewEvaluator_MissingQuery(t *testing.T) {
	_, err := opa.NewEvaluator(context.Background(), opa.Options{
		Module: `package helixgitpx`,
		Query:  "",
	})
	if err == nil {
		t.Fatal("expected error for empty query")
	}
}

func TestNewEvaluator_CompileError(t *testing.T) {
	_, err := opa.NewEvaluator(context.Background(), opa.Options{
		Module: `not valid rego {{{`,
		Query:  "data.x",
	})
	if err == nil {
		t.Fatal("expected compile error")
	}
}

func TestEval_ComplexPolicy(t *testing.T) {
	ev, err := opa.NewEvaluator(context.Background(), opa.Options{
		Module: `package helixgitpx
allow {
	input.role == "admin"
	input.org == input.target_org
}`,
		Query: "data.helixgitpx.allow",
	})
	if err != nil {
		t.Fatal(err)
	}
	got, err := ev.Eval(context.Background(), map[string]any{
		"role":        "admin",
		"org":         "acme",
		"target_org":  "acme",
	})
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := got.(bool); !b {
		t.Errorf("admin in same org should be allowed, got %v", got)
	}

	got2, err := ev.Eval(context.Background(), map[string]any{
		"role":        "admin",
		"org":         "acme",
		"target_org":  "other",
	})
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := got2.(bool); b {
		t.Errorf("admin in different org should be denied, got %v", got2)
	}
}
