// Package opa wraps github.com/open-policy-agent/opa/rego for in-process
// policy evaluation. M1 ships a minimal surface; M3 adds real policies.
package opa

import (
	"context"
	"errors"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
)

// Evaluator evaluates a compiled Rego query against structured input.
type Evaluator struct {
	query rego.PreparedEvalQuery
}

// Options configures NewEvaluator.
type Options struct {
	Module string // Rego source
	Query  string // e.g. "data.helixgitpx.allow"
}

// NewEvaluator compiles the given module and query.
func NewEvaluator(ctx context.Context, opts Options) (*Evaluator, error) {
	if opts.Query == "" {
		return nil, errors.New("opa: Query is required")
	}
	q, err := rego.New(
		rego.Query(opts.Query),
		rego.Module("policy.rego", opts.Module),
	).PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("opa: compile: %w", err)
	}
	return &Evaluator{query: q}, nil
}

// Eval runs the query with input and returns the first defined result or false.
func (e *Evaluator) Eval(ctx context.Context, input any) (any, error) {
	rs, err := e.query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return nil, err
	}
	if len(rs) == 0 || len(rs[0].Expressions) == 0 {
		return false, nil
	}
	return rs[0].Expressions[0].Value, nil
}
