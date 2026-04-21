// Package handler implements the HTTP surface opa-bundle-server exposes to
// in-cluster OPA agents via the `httpbundle` plugin, plus admin JSON endpoints.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/helixgitpx/helixgitpx/services/opa-bundle-server/internal/domain"
	"github.com/helixgitpx/helixgitpx/services/opa-bundle-server/internal/store"
)

type Handler struct {
	Store *store.Store
}

// Routes returns a ServeMux configured for the OPA bundle server surface.
func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /bundles/active", h.getActive)
	mux.HandleFunc("GET /bundles", h.list)
	mux.HandleFunc("GET /bundles/{id}", h.getOne)
	mux.HandleFunc("GET /healthz", h.healthz)
	return mux
}

func (h *Handler) getActive(w http.ResponseWriter, r *http.Request) {
	meta, content, err := h.Store.Active()
	if err != nil {
		http.Error(w, "no active bundle", http.StatusServiceUnavailable)
		return
	}
	if etag := r.Header.Get("If-None-Match"); etag != "" && etag == meta.ETag() {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("ETag", meta.ETag())
	w.Header().Set("Cache-Control", "no-transform, max-age=30")
	w.Header().Set("Content-Type", "application/vnd.openpolicyagent.bundles")
	w.Header().Set("X-Bundle-Id", meta.ID)
	w.Header().Set("X-Bundle-Version", meta.Version)
	w.Header().Set("X-Bundle-GitRev", meta.GitRev)
	_, _ = w.Write(content)
}

func (h *Handler) getOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	meta, content, err := h.Store.Get(id)
	if err != nil {
		if errors.Is(err, domain.ErrUnknownBundleID) {
			http.Error(w, "bundle not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("ETag", meta.ETag())
	w.Header().Set("Content-Type", "application/vnd.openpolicyagent.bundles")
	_, _ = w.Write(content)
}

func (h *Handler) list(w http.ResponseWriter, _ *http.Request) {
	bundles := h.Store.List()
	out := make([]map[string]any, 0, len(bundles))
	for _, b := range bundles {
		out = append(out, map[string]any{
			"id":        b.ID,
			"version":   b.Version,
			"git_rev":   b.GitRev,
			"active":    b.Active(),
			"size":      b.SizeBytes,
			"sha256":    b.Hash(),
			"signed":    b.Signed,
			"etag":      b.ETag(),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"bundles": out})
}

func (h *Handler) healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
