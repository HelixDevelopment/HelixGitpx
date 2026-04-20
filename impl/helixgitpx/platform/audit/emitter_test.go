package audit_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/helixgitpx/platform/audit"
)

func TestEvent_JSONEncoding(t *testing.T) {
	e := audit.Event{
		Action:      "org.create",
		Target:      "acme",
		ActorUserID: "user-1",
		Details:     map[string]any{"name": "Acme Inc"},
	}
	b, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var round audit.Event
	if err := json.Unmarshal(b, &round); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if round.Action != "org.create" || round.Target != "acme" {
		t.Errorf("roundtrip mismatch: %+v", round)
	}
	if round.At.IsZero() {
		t.Errorf("At should be set on Marshal when zero-valued")
	}
}

func TestEvent_AtPreservedWhenSet(t *testing.T) {
	t0 := time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC)
	e := audit.Event{At: t0, Action: "a", Target: "t"}
	b, _ := json.Marshal(e)
	var round audit.Event
	_ = json.Unmarshal(b, &round)
	if !round.At.Equal(t0) {
		t.Errorf("At = %v, want %v", round.At, t0)
	}
}
