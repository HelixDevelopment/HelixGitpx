package consumer

import (
	"encoding/json"
	"testing"
)

func TestRawEvent_JSONDecode(t *testing.T) {
	payload := []byte(`{"at":"2026-04-20T10:00:00Z","action":"org.create","target":"acme","actor_user_id":"u1"}`)
	var ev rawEvent
	if err := json.Unmarshal(payload, &ev); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if ev.Action != "org.create" {
		t.Errorf("Action = %q", ev.Action)
	}
}
