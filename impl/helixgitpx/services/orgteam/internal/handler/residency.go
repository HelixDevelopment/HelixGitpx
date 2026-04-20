// Package handler exposes orgteam RPCs over HTTP/Connect.
// This file wires the SetOrgResidency RPC to the domain layer.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/helixgitpx/helixgitpx/services/orgteam/internal/domain"
)

// OrgRepository abstracts the persistence surface for residency changes.
// Production uses the Postgres-backed repo; tests install a fake.
type OrgRepository interface {
	OwnerOf(ctx context.Context, orgID string) (string, error)
	SetResidency(ctx context.Context, orgID string, residency domain.Residency) error
}

// ResidencyHandler implements the HTTP surface for POST /v1/orgs/{id}/residency.
type ResidencyHandler struct {
	Repo        OrgRepository
	ActorFromCtx func(context.Context) string
}

type setResidencyRequest struct {
	Residency string `json:"residency"`
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ServeHTTP handles POST requests carrying a JSON body with the target residency.
// The org ID is parsed from the URL path prefix (caller strips `/v1/orgs/` and
// passes the id as `r.PathValue("id")`).
func (h *ResidencyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "use POST")
		return
	}
	orgID := r.PathValue("id")
	if orgID == "" {
		writeError(w, http.StatusBadRequest, "missing_org_id", "org id is required")
		return
	}

	var req setResidencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	actor := ""
	if h.ActorFromCtx != nil {
		actor = h.ActorFromCtx(r.Context())
	}
	if actor == "" {
		writeError(w, http.StatusUnauthorized, "unauthenticated", "actor not present on context")
		return
	}

	owner, err := h.Repo.OwnerOf(r.Context(), orgID)
	if err != nil {
		writeError(w, http.StatusNotFound, "org_not_found", err.Error())
		return
	}

	if err := domain.SetOrgResidency(owner, actor, domain.Residency(req.Residency)); err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidResidency):
			writeError(w, http.StatusBadRequest, "invalid_residency", err.Error())
		default:
			writeError(w, http.StatusForbidden, "forbidden", err.Error())
		}
		return
	}

	if err := h.Repo.SetResidency(r.Context(), orgID, domain.Residency(req.Residency)); err != nil {
		writeError(w, http.StatusInternalServerError, "persist_failed", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"id":        orgID,
		"residency": req.Residency,
	})
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{Code: code, Message: msg})
}
