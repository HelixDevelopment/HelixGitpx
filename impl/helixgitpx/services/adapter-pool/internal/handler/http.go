package handler

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"

	"github.com/helixgitpx/helixgitpx/services/adapter-pool/internal/adapter"
)

type Registry struct {
	mu       sync.Mutex
	provider map[adapter.Provider]adapter.Adapter
}

func NewRegistry() *Registry {
	return &Registry{provider: make(map[adapter.Provider]adapter.Adapter)}
}

func (r *Registry) Register(p adapter.Provider, a adapter.Adapter) {
	r.mu.Lock()
	r.provider[p] = a
	r.mu.Unlock()
}

func (r *Registry) List() []adapter.Provider {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]adapter.Provider, 0, len(r.provider))
	for p := range r.provider {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func (r *Registry) Get(p adapter.Provider) (adapter.Adapter, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.provider[p]
	return a, ok
}

type Handler struct {
	Registry *Registry
}

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/adapters", h.list)
	mux.HandleFunc("GET /v1/adapters/{provider}/health", h.health)
	mux.HandleFunc("GET /healthz", h.serviceHealth)
	return mux
}

func (h *Handler) list(w http.ResponseWriter, _ *http.Request) {
	providers := h.Registry.List()
	names := make([]string, len(providers))
	for i, p := range providers {
		names[i] = string(p)
	}
	writeJSON(w, http.StatusOK, map[string]any{"providers": names})
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	p := adapter.Provider(r.PathValue("provider"))
	a, ok := h.Registry.Get(p)
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "provider not registered")
		return
	}
	info, err := a.GetRepo(r.Context(), adapter.Source{Provider: p})
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "unhealthy", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"provider":  string(p),
		"healthy":   true,
		"default_branch": info.Default,
	})
}

func (h *Handler) serviceHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, code int, kind, msg string) {
	writeJSON(w, code, map[string]string{"code": kind, "message": msg})
}
