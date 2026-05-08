package engines

import (
	"context"
	"testing"
)

func TestQuery_Fields(t *testing.T) {
	q := Query{
		Text:   "hello",
		OrgID:  "org1",
		Repos:  []string{"repo1", "repo2"},
		Limit:  10,
		Offset: 5,
	}
	if q.Text != "hello" {
		t.Errorf("Text = %q", q.Text)
	}
	if len(q.Repos) != 2 {
		t.Errorf("Repos len = %d, want 2", len(q.Repos))
	}
	if q.Limit != 10 {
		t.Errorf("Limit = %d, want 10", q.Limit)
	}
}

func TestHit_ZeroValue(t *testing.T) {
	var h Hit
	if h.ID != "" {
		t.Errorf("zero Hit.ID should be empty, got %q", h.ID)
	}
	if h.Score != 0 {
		t.Errorf("zero Hit.Score should be 0, got %f", h.Score)
	}
}

type stubEngine struct {
	name  string
	hits  []Hit
	calls int
}

func (s *stubEngine) Name() string { return s.name }
func (s *stubEngine) Search(_ context.Context, q Query) ([]Hit, error) {
	s.calls++
	return s.hits, nil
}

func TestStubEngine_ImplementsInterface(t *testing.T) {
	var _ Engine = (*stubEngine)(nil)
}

func TestStubEngine_SearchReturnsHits(t *testing.T) {
	e := &stubEngine{
		name: "test",
		hits: []Hit{{ID: "1", Title: "result", Score: 0.9}},
	}
	hits, err := e.Search(context.Background(), Query{Text: "query"})
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 {
		t.Fatalf("hits = %d, want 1", len(hits))
	}
	if hits[0].Title != "result" {
		t.Errorf("Title = %q, want result", hits[0].Title)
	}
	if e.calls != 1 {
		t.Errorf("calls = %d, want 1", e.calls)
	}
}

func TestQuery_EmptyRepos(t *testing.T) {
	q := Query{Text: "test", Repos: nil}
	if q.Repos != nil {
		t.Errorf("nil Repos should stay nil")
	}
}

func TestHit_AllFields(t *testing.T) {
	h := Hit{
		ID:      "h1",
		Kind:    "code",
		Title:   "main.go",
		Snippet: "func main()",
		URL:     "https://example.com",
		Score:   0.95,
	}
	if h.Kind != "code" {
		t.Errorf("Kind = %q", h.Kind)
	}
	if h.Snippet != "func main()" {
		t.Errorf("Snippet = %q", h.Snippet)
	}
}
