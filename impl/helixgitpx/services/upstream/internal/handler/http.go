// Package handler is the HTTP surface of upstream-service.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/helixgitpx/helixgitpx/services/upstream/internal/domain"
	"github.com/helixgitpx/helixgitpx/services/upstream/internal/memstore"
)

type Store interface {
	Create(ctx context.Context, in domain.BindingInput) (memstore.Binding, error)
	Get(ctx context.Context, id string) (memstore.Binding, error)
	ListByRepo(ctx context.Context, repoID string) []memstore.Binding
	Delete(ctx context.Context, id string) error
}

type Handler struct {
	Store Store
}

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/upstreams", h.create)
	mux.HandleFunc("GET /v1/upstreams", h.list)
	mux.HandleFunc("GET /v1/upstreams/{id}", h.get)
	mux.HandleFunc("DELETE /v1/upstreams/{id}", h.del)
	mux.HandleFunc("GET /healthz", h.health)
	return mux
}

type bindingIn struct {
	RepoID    string `json:"repo_id"`
	Provider  string `json:"provider"`
	URL       string `json:"url"`
	Direction string `json:"direction"`
}
type bindingOut struct {
	ID        string `json:"id"`
	RepoID    string `json:"repo_id"`
	Provider  string `json:"provider"`
	URL       string `json:"url"`
	Direction string `json:"direction"`
}

func directionFromString(s string) domain.Direction {
	switch s {
	case "read":
		return domain.DirectionReadOnly
	case "write":
		return domain.DirectionWrite
	case "bidirectional", "both":
		return domain.DirectionBidirectional
	}
	return domain.DirectionUnspecified
}

func directionToString(d domain.Direction) string {
	switch d {
	case domain.DirectionReadOnly:
		return "read"
	case domain.DirectionWrite:
		return "write"
	case domain.DirectionBidirectional:
		return "bidirectional"
	}
	return ""
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in bindingIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	input := domain.BindingInput{
		RepoID:    in.RepoID,
		Provider:  in.Provider,
		RawURL:    in.URL,
		Direction: directionFromString(in.Direction),
	}
	if err := domain.Validate(input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid", err.Error())
		return
	}
	b, err := h.Store.Create(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "create_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, bindingOut{
		ID: b.ID, RepoID: b.RepoID, Provider: b.Provider, URL: b.URL,
		Direction: directionToString(b.Direction),
	})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	repoID := r.URL.Query().Get("repo_id")
	if repoID == "" {
		writeError(w, http.StatusBadRequest, "missing_repo_id", "repo_id query param required")
		return
	}
	bindings := h.Store.ListByRepo(r.Context(), repoID)
	out := make([]bindingOut, 0, len(bindings))
	for _, b := range bindings {
		out = append(out, bindingOut{
			ID: b.ID, RepoID: b.RepoID, Provider: b.Provider, URL: b.URL,
			Direction: directionToString(b.Direction),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"bindings": out})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := h.Store.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, memstore.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not_found", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "get_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, bindingOut{
		ID: b.ID, RepoID: b.RepoID, Provider: b.Provider, URL: b.URL,
		Direction: directionToString(b.Direction),
	})
}

func (h *Handler) del(w http.ResponseWriter, r *http.Request) {
	if err := h.Store.Delete(r.Context(), r.PathValue("id")); err != nil {
		if errors.Is(err, memstore.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not_found", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "delete_failed", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
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
