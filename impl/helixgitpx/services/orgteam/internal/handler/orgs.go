// Package handler — org CRUD endpoints (list/create/get).
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/helixgitpx/helixgitpx/services/orgteam/internal/domain"
)

// OrgStore is the wider persistence surface required for the org endpoints.
// ResidencyHandler's OrgRepository is a subset; concrete memstore.Store
// satisfies both.
type OrgStore interface {
	OrgRepository
	CreateOrg(ctx context.Context, slug, name, ownerID string, residency domain.Residency) (domain.Org, error)
	GetOrg(ctx context.Context, orgID string) (domain.Org, domain.Residency, error)
	ListOrgs(ctx context.Context) []domain.Org
}

// ErrDuplicateSlug flags the 409 from Store.CreateOrg.
var ErrDuplicateSlug = errors.New("orgteam: slug already exists")

type OrgHandler struct {
	Store        OrgStore
	ActorFromCtx func(context.Context) string
}

type orgIn struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Residency string `json:"residency"`
}
type orgOut struct {
	ID        string `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Residency string `json:"residency,omitempty"`
}

// RegisterOrgRoutes attaches list/create/get to an existing mux.
func (h *OrgHandler) RegisterOrgRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/orgs", h.list)
	mux.HandleFunc("POST /v1/orgs", h.create)
	mux.HandleFunc("GET /v1/orgs/{id}", h.get)
}

func (h *OrgHandler) list(w http.ResponseWriter, r *http.Request) {
	orgs := h.Store.ListOrgs(r.Context())
	out := make([]orgOut, 0, len(orgs))
	for _, o := range orgs {
		out = append(out, orgOut{ID: o.ID, Slug: o.Slug, Name: o.Name})
	}
	writeJSON(w, http.StatusOK, map[string]any{"orgs": out})
}

func (h *OrgHandler) create(w http.ResponseWriter, r *http.Request) {
	var in orgIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if strings.TrimSpace(in.Slug) == "" || strings.TrimSpace(in.Name) == "" {
		writeError(w, http.StatusBadRequest, "missing_fields", "slug and name are required")
		return
	}
	residency := domain.Residency(strings.ToUpper(strings.TrimSpace(in.Residency)))
	if residency == "" {
		residency = domain.ResidencyEU
	}
	if !residency.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_residency", string(residency)+" is not a supported zone")
		return
	}

	actor := ""
	if h.ActorFromCtx != nil {
		actor = h.ActorFromCtx(r.Context())
	}
	if actor == "" {
		writeError(w, http.StatusUnauthorized, "unauthenticated", "actor required")
		return
	}

	org, err := h.Store.CreateOrg(r.Context(), in.Slug, in.Name, actor, residency)
	if err != nil {
		if errors.Is(err, ErrDuplicateSlug) || strings.Contains(err.Error(), "already exists") {
			writeError(w, http.StatusConflict, "duplicate_slug", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "create_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, orgOut{ID: org.ID, Slug: org.Slug, Name: org.Name, Residency: string(residency)})
}

func (h *OrgHandler) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	org, residency, err := h.Store.GetOrg(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, orgOut{ID: org.ID, Slug: org.Slug, Name: org.Name, Residency: string(residency)})
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}
