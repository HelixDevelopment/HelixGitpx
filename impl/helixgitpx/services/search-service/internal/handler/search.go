// Package handler is the HTTP surface of search-service.
// It accepts a query, fans out to each registered engine in parallel,
// applies RRF fusion, and returns top-K.
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/helixgitpx/helixgitpx/services/search-service/internal/domain"
	"github.com/helixgitpx/helixgitpx/services/search-service/internal/engines"
)

type Handler struct {
	Engines []engines.Engine
	Timeout time.Duration
}

// Routes returns the registered HTTP mux.
func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/search", h.search)
	mux.HandleFunc("GET /healthz", h.healthz)
	return mux
}

type response struct {
	Hits      []responseHit `json:"hits"`
	ElapsedMs int64         `json:"elapsed_ms"`
	Engines   []string      `json:"engines"`
}

type responseHit struct {
	ID        string             `json:"id"`
	Score     float64            `json:"score"`
	PerEngine map[string]float64 `json:"per_engine,omitempty"`
}

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	org := r.URL.Query().Get("org_id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}
	timeout := h.Timeout
	if timeout == 0 {
		timeout = 2 * time.Second
	}
	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	start := time.Now()
	rankings := fanOut(ctx, h.Engines, engines.Query{Text: q, OrgID: org, Limit: limit})
	fused := domain.TopK(domain.Fuse(rankings), limit)

	engineNames := make([]string, 0, len(h.Engines))
	for _, e := range h.Engines {
		engineNames = append(engineNames, e.Name())
	}

	out := response{
		Hits:      make([]responseHit, 0, len(fused)),
		ElapsedMs: time.Since(start).Milliseconds(),
		Engines:   engineNames,
	}
	for _, f := range fused {
		out.Hits = append(out.Hits, responseHit{ID: f.ID, Score: f.Score, PerEngine: f.PerEngine})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func (h *Handler) healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func fanOut(ctx context.Context, es []engines.Engine, q engines.Query) []domain.Ranking {
	var wg sync.WaitGroup
	results := make([]domain.Ranking, len(es))
	for i, e := range es {
		wg.Add(1)
		go func(i int, e engines.Engine) {
			defer wg.Done()
			hits, err := e.Search(ctx, q)
			if err != nil {
				return
			}
			ids := make([]string, 0, len(hits))
			for _, h := range hits {
				ids = append(ids, h.ID)
			}
			results[i] = domain.Ranking{Engine: e.Name(), Hits: ids, Weight: 1.0}
		}(i, e)
	}
	wg.Wait()
	return results
}
