package handler

import (
	"encoding/json"
	"net/http"

	"github.com/helixgitpx/helixgitpx/services/collab-service/internal/domain"
)

type Handler struct {
	Limits domain.Limits
}

func NewHandler() *Handler {
	return &Handler{Limits: domain.DefaultLimits()}
}

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/docs/validate", h.validateDoc)
	mux.HandleFunc("POST /v1/docs/snapshot-check", h.snapshotCheck)
	mux.HandleFunc("POST /v1/docs/participant-check", h.participantCheck)
	mux.HandleFunc("GET /healthz", h.health)
	return mux
}

type validateDocIn struct {
	DocID  string `json:"doc_id"`
	ActorID string `json:"actor_id"`
}

type validateDocOut struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

func (h *Handler) validateDoc(w http.ResponseWriter, r *http.Request) {
	var in validateDocIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if err := domain.ValidateOpenDoc(in.DocID, in.ActorID); err != nil {
		writeJSON(w, http.StatusOK, validateDocOut{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, validateDocOut{Valid: true})
}

type snapshotCheckIn struct {
	Snapshot []byte `json:"snapshot"`
}

type snapshotCheckOut struct {
	Allowed bool   `json:"allowed"`
	Error   string `json:"error,omitempty"`
}

func (h *Handler) snapshotCheck(w http.ResponseWriter, r *http.Request) {
	var in snapshotCheckIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if err := domain.SnapshotSizeAllowed(h.Limits, in.Snapshot); err != nil {
		writeJSON(w, http.StatusOK, snapshotCheckOut{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, snapshotCheckOut{Allowed: true})
}

type participantCheckIn struct {
	CurrentParticipants int `json:"current_participants"`
}

type participantCheckOut struct {
	Allowed bool   `json:"allowed"`
	Error   string `json:"error,omitempty"`
}

func (h *Handler) participantCheck(w http.ResponseWriter, r *http.Request) {
	var in participantCheckIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if err := domain.AddParticipantAllowed(h.Limits, in.CurrentParticipants); err != nil {
		writeJSON(w, http.StatusOK, participantCheckOut{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, participantCheckOut{Allowed: true})
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
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
