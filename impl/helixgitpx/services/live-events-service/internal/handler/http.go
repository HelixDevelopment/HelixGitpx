package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/helixgitpx/helixgitpx/services/live-events-service/internal/domain"
)

type Handler struct{}

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/tokens/encode", h.encodeToken)
	mux.HandleFunc("POST /v1/tokens/decode", h.decodeToken)
	mux.HandleFunc("POST /v1/events/match", h.matchEvent)
	mux.HandleFunc("GET /healthz", h.health)
	return mux
}

type encodeIn struct {
	Offset    int64 `json:"offset"`
	Timestamp int64 `json:"timestamp"`
}

type encodeOut struct {
	Token string `json:"token"`
}

func (h *Handler) encodeToken(w http.ResponseWriter, r *http.Request) {
	var in encodeIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	tok := domain.ResumeToken{
		Offset:    in.Offset,
		Timestamp: time.Unix(in.Timestamp, 0).UTC(),
	}
	writeJSON(w, http.StatusOK, encodeOut{Token: tok.Encode()})
}

type decodeIn struct {
	Token     string `json:"token"`
	Retention int64  `json:"retention_seconds"`
}

type decodeOut struct {
	Offset    int64 `json:"offset"`
	Timestamp int64 `json:"timestamp"`
	Error     string `json:"error,omitempty"`
}

func (h *Handler) decodeToken(w http.ResponseWriter, r *http.Request) {
	var in decodeIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	retention := time.Duration(in.Retention) * time.Second
	tok, err := domain.DecodeResumeToken(in.Token, retention, time.Now().UTC())
	if err != nil {
		writeJSON(w, http.StatusOK, decodeOut{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, decodeOut{
		Offset:    tok.Offset,
		Timestamp: tok.Timestamp.Unix(),
	})
}

type matchIn struct {
	SubRepos []string `json:"sub_repos"`
	SubTypes []string `json:"sub_types"`
	RepoID   string   `json:"repo_id"`
	EventType string  `json:"event_type"`
}

type matchOut struct {
	Match bool `json:"match"`
}

func (h *Handler) matchEvent(w http.ResponseWriter, r *http.Request) {
	var in matchIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, matchOut{
		Match: domain.Matches(in.SubRepos, in.SubTypes, in.RepoID, in.EventType),
	})
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
