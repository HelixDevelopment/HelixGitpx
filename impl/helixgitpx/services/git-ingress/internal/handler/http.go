package handler

import (
	"encoding/json"
	"net/http"

	"github.com/helixgitpx/helixgitpx/services/git-ingress/internal/domain"
)

type Handler struct{}

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/pushes/validate", h.validatePush)
	mux.HandleFunc("GET /healthz", h.health)
	return mux
}

type validateIn struct {
	Ref             string `json:"ref"`
	RepoID          string `json:"repo_id"`
	SizeBytes       int64  `json:"size_bytes"`
	PushesLastMinute int   `json:"pushes_last_minute"`
	PushLimit       int    `json:"push_limit"`
	MaxBytesPerPush int64  `json:"max_bytes_per_push"`
}

type validateOut struct {
	Valid     bool   `json:"valid"`
	Protected bool   `json:"protected"`
	Error     string `json:"error,omitempty"`
}

func (h *Handler) validatePush(w http.ResponseWriter, r *http.Request) {
	var in validateIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	refErr := domain.ValidateRef(in.Ref)
	if refErr != nil {
		writeJSON(w, http.StatusOK, validateOut{Valid: false, Error: refErr.Error()})
		return
	}

	pushErr := domain.AllowPush(domain.AllowPushInput{
		RepoID:           in.RepoID,
		SizeBytes:        in.SizeBytes,
		PushesLastMinute: in.PushesLastMinute,
		PushLimit:        in.PushLimit,
		MaxBytesPerPush:  in.MaxBytesPerPush,
	})
	if pushErr != nil {
		writeJSON(w, http.StatusOK, validateOut{
			Valid:     false,
			Protected: domain.IsProtected(in.Ref),
			Error:     pushErr.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, validateOut{
		Valid:     true,
		Protected: domain.IsProtected(in.Ref),
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
