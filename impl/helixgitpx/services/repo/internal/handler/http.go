// Package handler owns the HTTP surface of repo-service: list, create,
// get, delete repositories + manage branch protections.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/helixgitpx/helixgitpx/services/repo/internal/domain"
)

// Store is the persistence port. Production wires it to Postgres; tests
// install an in-memory fake.
type Store interface {
	Create(ctx context.Context, r domain.Repo) (domain.Repo, error)
	Get(ctx context.Context, id string) (domain.Repo, error)
	List(ctx context.Context, orgID string) ([]domain.Repo, error)
	Delete(ctx context.Context, id string) error

	AddProtection(ctx context.Context, p domain.Protection) error
	ListProtections(ctx context.Context, repoID string) ([]domain.Protection, error)
}

// ErrNotFound signals a missing repo.
var ErrNotFound = errors.New("repo: not found")

type Handler struct {
	Store Store
}

// Routes returns the service's ServeMux.
func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/repos", h.list)
	mux.HandleFunc("POST /v1/repos", h.create)
	mux.HandleFunc("GET /v1/repos/{id}", h.get)
	mux.HandleFunc("DELETE /v1/repos/{id}", h.del)
	mux.HandleFunc("POST /v1/repos/{id}/protections", h.addProtection)
	mux.HandleFunc("GET /v1/repos/{id}/protections", h.listProtections)
	mux.HandleFunc("GET /healthz", h.healthz)
	return mux
}

type repoIn struct {
	OrgID         string `json:"org_id"`
	Slug          string `json:"slug"`
	DefaultBranch string `json:"default_branch"`
}
type repoOut struct {
	ID            string `json:"id"`
	OrgID         string `json:"org_id"`
	Slug          string `json:"slug"`
	DefaultBranch string `json:"default_branch"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in repoIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if strings.TrimSpace(in.OrgID) == "" || strings.TrimSpace(in.Slug) == "" {
		writeErr(w, http.StatusBadRequest, "missing_fields", "org_id and slug are required")
		return
	}
	if in.DefaultBranch == "" {
		in.DefaultBranch = "main"
	}
	repo, err := h.Store.Create(r.Context(), domain.Repo{
		OrgID:         in.OrgID,
		Slug:          in.Slug,
		DefaultBranch: in.DefaultBranch,
	})
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "create_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, repoOut{ID: repo.ID, OrgID: repo.OrgID, Slug: repo.Slug, DefaultBranch: repo.DefaultBranch})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	orgID := r.URL.Query().Get("org_id")
	repos, err := h.Store.List(r.Context(), orgID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}
	out := make([]repoOut, 0, len(repos))
	for _, r := range repos {
		out = append(out, repoOut{ID: r.ID, OrgID: r.OrgID, Slug: r.Slug, DefaultBranch: r.DefaultBranch})
	}
	writeJSON(w, http.StatusOK, map[string]any{"repos": out})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	repo, err := h.Store.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeErr(w, http.StatusNotFound, "not_found", err.Error())
			return
		}
		writeErr(w, http.StatusInternalServerError, "get_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, repoOut{ID: repo.ID, OrgID: repo.OrgID, Slug: repo.Slug, DefaultBranch: repo.DefaultBranch})
}

func (h *Handler) del(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.Store.Delete(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeErr(w, http.StatusNotFound, "not_found", err.Error())
			return
		}
		writeErr(w, http.StatusInternalServerError, "delete_failed", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type protectionIn struct {
	Pattern           string `json:"pattern"`
	RequireSigned     bool   `json:"require_signed"`
	RequiredReviewers int    `json:"required_reviewers"`
}

func (h *Handler) addProtection(w http.ResponseWriter, r *http.Request) {
	repoID := r.PathValue("id")
	var in protectionIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if strings.TrimSpace(in.Pattern) == "" {
		writeErr(w, http.StatusBadRequest, "missing_pattern", "pattern is required")
		return
	}
	if in.RequiredReviewers < 0 {
		writeErr(w, http.StatusBadRequest, "bad_reviewers", "required_reviewers must be non-negative")
		return
	}
	if err := h.Store.AddProtection(r.Context(), domain.Protection{
		RepoID:            repoID,
		Pattern:           in.Pattern,
		RequireSigned:     in.RequireSigned,
		RequiredReviewers: in.RequiredReviewers,
	}); err != nil {
		writeErr(w, http.StatusInternalServerError, "protection_failed", err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) listProtections(w http.ResponseWriter, r *http.Request) {
	list, err := h.Store.ListProtections(r.Context(), r.PathValue("id"))
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"protections": list})
}

func (h *Handler) healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func writeErr(w http.ResponseWriter, code int, kind, msg string) {
	writeJSON(w, code, map[string]string{"code": kind, "message": msg})
}
