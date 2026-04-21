// Package handler is the HTTP surface of billing-service.
package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/helixgitpx/helixgitpx/services/billing-service/internal/usecase"
)

type Handler struct {
	UseCases *usecase.UseCases
}

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/subscriptions/{id}/upgrade", h.upgrade)
	mux.HandleFunc("POST /v1/subscriptions/{id}/cancel", h.cancel)
	mux.HandleFunc("GET /v1/plans", h.listPlans)
	mux.HandleFunc("GET /healthz", h.health)
	return mux
}

type upgradeIn struct {
	OrgID string `json:"org_id"`
	Plan  string `json:"plan"`
}

func (h *Handler) upgrade(w http.ResponseWriter, r *http.Request) {
	subID := r.PathValue("id")
	var in upgradeIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	sub, err := h.UseCases.UpgradePlan(context.Background(), in.OrgID, subID, usecase.PlanName(in.Plan))
	if err != nil {
		writeError(w, http.StatusBadRequest, "upgrade_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"subscription_id": sub.ExternalID,
		"plan":            sub.Plan,
		"status":          sub.Status,
	})
}

func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
	subID := r.PathValue("id")
	if err := h.UseCases.CancelPlan(context.Background(), subID); err != nil {
		writeError(w, http.StatusInternalServerError, "cancel_failed", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listPlans(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"plans": []string{
			string(usecase.PlanFree),
			string(usecase.PlanTeam),
			string(usecase.PlanScale),
			string(usecase.PlanEnt),
		},
	})
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
