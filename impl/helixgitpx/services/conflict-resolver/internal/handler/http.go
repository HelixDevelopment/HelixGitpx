package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/helixgitpx/helixgitpx/services/conflict-resolver/internal/domain"
)

type Store struct {
	mu   sync.Mutex
	next uint64
	byID map[string]entry
}

type entry struct {
	Conflict domainEntry
	Status   domain.Status
	Rationale string
}

type domainEntry struct {
	ID             string
	Kind           domain.Kind
	RepoID         string
	UpstreamRef    string
	UpstreamLabel  string
	LocalRef       string
	LocalLabel     string
	RenamedTo      string
	MetaKey        string
}

func NewStore() *Store {
	return &Store{byID: make(map[string]entry)}
}

func (s *Store) Create(de domainEntry) (string, domain.Status) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := fmt.Sprintf("conf-%d", atomic.AddUint64(&s.next, 1))
	de.ID = id
	s.byID[id] = entry{Conflict: de, Status: domain.StatusOpen}
	return id, domain.StatusOpen
}

func (s *Store) Get(id string) (entry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.byID[id]
	return e, ok
}

func (s *Store) SetStatus(id string, st domain.Status, rationale string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.byID[id]
	if !ok {
		return fmt.Errorf("not found: %s", id)
	}
	if err := domain.Transition(e.Status, st); err != nil {
		return err
	}
	e.Status = st
	e.Rationale = rationale
	s.byID[id] = e
	return nil
}

func (s *Store) Propose(id, rationale string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.byID[id]
	if !ok {
		return fmt.Errorf("not found: %s", id)
	}
	if err := domain.Transition(e.Status, domain.StatusProposed); err != nil {
		return err
	}
	if err := domain.ValidateRationale("human", rationale); err != nil {
		return err
	}
	e.Status = domain.StatusProposed
	e.Rationale = rationale
	s.byID[id] = e
	return nil
}

type Handler struct {
	Store *Store
}

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/conflicts", h.create)
	mux.HandleFunc("GET /v1/conflicts/{id}", h.get)
	mux.HandleFunc("POST /v1/conflicts/{id}/propose", h.propose)
	mux.HandleFunc("POST /v1/conflicts/{id}/resolve", h.resolve)
	mux.HandleFunc("POST /v1/conflicts/{id}/reject", h.reject)
	mux.HandleFunc("GET /healthz", h.health)
	return mux
}

type createIn struct {
	RepoID        string `json:"repo_id"`
	RefsDiverge   bool   `json:"refs_diverge"`
	LabelsDiffer  bool   `json:"labels_differ"`
	RenameCollision bool `json:"rename_collision"`
	MetaDrift     bool   `json:"meta_drift"`
}

type conflictOut struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Status string `json:"status"`
	RepoID string `json:"repo_id"`
}

func kindToString(k domain.Kind) string {
	return k.String()
}

func statusToString(s domain.Status) string {
	switch s {
	case domain.StatusOpen:
		return "open"
	case domain.StatusProposed:
		return "proposed"
	case domain.StatusResolved:
		return "resolved"
	case domain.StatusRejected:
		return "rejected"
	default:
		return "unspecified"
	}
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in createIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	kind := domain.Classify(in.RefsDiverge, in.LabelsDiffer, in.RenameCollision, in.MetaDrift)
	if kind == domain.KindUnspecified {
		writeError(w, http.StatusBadRequest, "no_conflict", "no conflict signal provided")
		return
	}
	de := domainEntry{
		Kind:   kind,
		RepoID: in.RepoID,
	}
	id, st := h.Store.Create(de)
	writeJSON(w, http.StatusCreated, conflictOut{
		ID: id, Kind: kindToString(kind), Status: statusToString(st), RepoID: in.RepoID,
	})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, ok := h.Store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "conflict not found")
		return
	}
	writeJSON(w, http.StatusOK, conflictOut{
		ID: e.Conflict.ID, Kind: kindToString(e.Conflict.Kind),
		Status: statusToString(e.Status), RepoID: e.Conflict.RepoID,
	})
}

type proposeIn struct {
	Rationale string `json:"rationale"`
}

func (h *Handler) propose(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var in proposeIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if err := h.Store.Propose(id, in.Rationale); err != nil {
		writeError(w, http.StatusBadRequest, "propose_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "proposed"})
}

type resolveIn struct {
	Rationale string `json:"rationale"`
}

func (h *Handler) resolve(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var in resolveIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if err := h.Store.SetStatus(id, domain.StatusResolved, in.Rationale); err != nil {
		writeError(w, http.StatusBadRequest, "resolve_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "resolved"})
}

func (h *Handler) reject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var in resolveIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if err := h.Store.SetStatus(id, domain.StatusRejected, in.Rationale); err != nil {
		writeError(w, http.StatusBadRequest, "reject_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "rejected"})
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
