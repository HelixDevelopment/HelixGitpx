// Package engines abstracts the three backends search-service fuses:
// Meilisearch (keyword), Qdrant (vector), Zoekt (code regex).
package engines

import "context"

// Query is the normalized input passed to each engine.
type Query struct {
	Text   string
	OrgID  string
	Repos  []string
	Limit  int
	Offset int
}

// Engine is the contract every backend must satisfy.
type Engine interface {
	Name() string
	Search(ctx context.Context, q Query) ([]Hit, error)
}

// Hit is a single result with provenance back to the engine that produced it.
type Hit struct {
	ID      string
	Kind    string
	Title   string
	Snippet string
	URL     string
	Score   float64
}
