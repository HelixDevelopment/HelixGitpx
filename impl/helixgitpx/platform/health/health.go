// Package health provides /livez, /readyz, /healthz HTTP handlers and a
// probe registry. Register a named probe with Register; Ready returns 503
// when any probe fails.
package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
)

// Probe returns nil when the dependency is healthy.
type Probe func(context.Context) error

// Handler serves the three liveness/readiness endpoints.
type Handler struct {
	mu     sync.RWMutex
	probes map[string]Probe
}

// New builds an empty Handler.
func New() *Handler {
	return &Handler{probes: make(map[string]Probe)}
}

// Register attaches a probe under name; re-registering replaces the prior probe.
func (h *Handler) Register(name string, p Probe) {
	h.mu.Lock()
	h.probes[name] = p
	h.mu.Unlock()
}

// Routes registers handler funcs on the given mux.
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", h.Live)
	mux.HandleFunc("/livez", h.Live)
	mux.HandleFunc("/readyz", h.Ready)
}

// Live always returns 200; process is up.
func (h *Handler) Live(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}

// Ready runs every probe with the request context and returns 503 if any fail.
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	probes := make(map[string]Probe, len(h.probes))
	for k, v := range h.probes {
		probes[k] = v
	}
	h.mu.RUnlock()

	results := make(map[string]any, len(probes))
	allOK := true
	for name, p := range probes {
		if err := p(r.Context()); err != nil {
			results[name] = err.Error()
			allOK = false
		} else {
			results[name] = "ok"
		}
	}

	status := http.StatusOK
	body := map[string]any{"status": "ok", "checks": results}
	if !allOK {
		status = http.StatusServiceUnavailable
		body["status"] = "unavailable"
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
