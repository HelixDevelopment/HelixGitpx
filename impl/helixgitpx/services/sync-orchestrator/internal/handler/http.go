package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/helixgitpx/helixgitpx/services/sync-orchestrator/internal/domain"
)

type Handler struct{}

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/retry/classify", h.classify)
	mux.HandleFunc("POST /v1/retry/backoff", h.backoff)
	mux.HandleFunc("POST /v1/retry/decision", h.decision)
	mux.HandleFunc("GET /healthz", h.health)
	return mux
}

type classifyIn struct {
	HTTPStatus int    `json:"http_status"`
	ErrMessage string `json:"err_message"`
	Permanent  bool   `json:"permanent"`
}

type classifyOut struct {
	Kind string `json:"kind"`
}

func kindToString(k domain.ErrorKind) string {
	switch k {
	case domain.KindTransient:
		return "transient"
	case domain.KindRateLimit:
		return "rate_limit"
	case domain.KindAuthFailed:
		return "auth_failed"
	case domain.KindClientError:
		return "client_error"
	case domain.KindPermanent:
		return "permanent"
	default:
		return "unknown"
	}
}

func (h *Handler) classify(w http.ResponseWriter, r *http.Request) {
	var in classifyIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	var err error
	if in.Permanent {
		err = domain.ErrPermanentSentinel
	} else if in.ErrMessage != "" {
		err = newErr(in.ErrMessage)
	}
	kind := domain.Classify(in.HTTPStatus, err)
	writeJSON(w, http.StatusOK, classifyOut{Kind: kindToString(kind)})
}

type backoffIn struct {
	Attempt    int   `json:"attempt"`
	BaseMs     int64 `json:"base_ms"`
	MaxMs      int64 `json:"max_ms"`
}

type backoffOut struct {
	DelayMs int64 `json:"delay_ms"`
}

func (h *Handler) backoff(w http.ResponseWriter, r *http.Request) {
	var in backoffIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	d := domain.Backoff(in.Attempt, time.Duration(in.BaseMs)*time.Millisecond, time.Duration(in.MaxMs)*time.Millisecond)
	writeJSON(w, http.StatusOK, backoffOut{DelayMs: d.Milliseconds()})
}

type decisionIn struct {
	HTTPStatus  int  `json:"http_status"`
	Permanent   bool `json:"permanent"`
	Attempt     int  `json:"attempt"`
	MaxAttempts int  `json:"max_attempts"`
}

type decisionOut struct {
	Kind       string `json:"kind"`
	ShouldRetry bool  `json:"should_retry"`
	GoesToDLQ  bool   `json:"goes_to_dlq"`
}

func (h *Handler) decision(w http.ResponseWriter, r *http.Request) {
	var in decisionIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	var err error
	if in.Permanent {
		err = domain.ErrPermanentSentinel
	}
	kind := domain.Classify(in.HTTPStatus, err)
	writeJSON(w, http.StatusOK, decisionOut{
		Kind:        kindToString(kind),
		ShouldRetry: domain.ShouldRetry(kind, in.Attempt, in.MaxAttempts),
		GoesToDLQ:   domain.GoesToDLQ(kind, in.Attempt, in.MaxAttempts),
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

type simpleErr struct{ msg string }

func (e *simpleErr) Error() string { return e.msg }

func newErr(msg string) error { return &simpleErr{msg: msg} }
