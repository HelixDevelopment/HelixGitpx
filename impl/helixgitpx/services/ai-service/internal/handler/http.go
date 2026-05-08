package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/helixgitpx/helixgitpx/services/ai-service/internal/usecase"
)

type Handler struct {
	UseCases *usecase.UseCases
}

func NewHandler(uc *usecase.UseCases) *Handler {
	return &Handler{UseCases: uc}
}

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/ai/summarize", h.summarize)
	mux.HandleFunc("POST /v1/ai/conflict", h.conflict)
	mux.HandleFunc("POST /v1/ai/labels", h.labels)
	mux.HandleFunc("POST /v1/ai/chatops", h.chatops)
	mux.HandleFunc("GET /healthz", h.health)
	return mux
}

type summarizeIn struct {
	Content string `json:"content"`
}

type promptOut struct {
	Result string `json:"result"`
}

func (h *Handler) summarize(w http.ResponseWriter, r *http.Request) {
	var in summarizeIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	result, err := h.UseCases.Summarize(r.Context(), in.Content)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "summarize_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, promptOut{Result: result})
}

type conflictIn struct {
	Diff string `json:"diff"`
}

func (h *Handler) conflict(w http.ResponseWriter, r *http.Request) {
	var in conflictIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	result, err := h.UseCases.ProposeConflict(r.Context(), in.Diff)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "conflict_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, promptOut{Result: result})
}

type labelsIn struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (h *Handler) labels(w http.ResponseWriter, r *http.Request) {
	var in labelsIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	result, err := h.UseCases.SuggestLabel(r.Context(), in.Title, in.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "labels_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, promptOut{Result: result})
}

type chatopsIn struct {
	Prompt string `json:"prompt"`
}

func (h *Handler) chatops(w http.ResponseWriter, r *http.Request) {
	var in chatopsIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	result, err := h.UseCases.ChatOps(r.Context(), in.Prompt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "chatops_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, promptOut{Result: result})
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

var _ = context.TODO
