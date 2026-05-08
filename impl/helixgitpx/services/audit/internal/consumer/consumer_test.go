package consumer

import (
	"encoding/json"
	"testing"
	"time"
)

func TestRawEvent_JSONDecode(t *testing.T) {
	payload := []byte(`{"at":"2026-04-20T10:00:00Z","actor_user_id":"u1","actor_ip":"10.0.0.1","action":"org.create","target":"acme","details":{"region":"EU"}}`)
	var ev rawEvent
	if err := json.Unmarshal(payload, &ev); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if ev.Action != "org.create" {
		t.Errorf("Action = %q, want %q", ev.Action, "org.create")
	}
	if ev.ActorUserID != "u1" {
		t.Errorf("ActorUserID = %q, want %q", ev.ActorUserID, "u1")
	}
	if ev.ActorIP != "10.0.0.1" {
		t.Errorf("ActorIP = %q, want %q", ev.ActorIP, "10.0.0.1")
	}
	if ev.Target != "acme" {
		t.Errorf("Target = %q, want %q", ev.Target, "acme")
	}
	wantTime := time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC)
	if !ev.At.Equal(wantTime) {
		t.Errorf("At = %v, want %v", ev.At, wantTime)
	}
	if ev.Details["region"] != "EU" {
		t.Errorf("Details[region] = %v, want %q", ev.Details["region"], "EU")
	}
}

func TestRawEvent_JSONDecode_MissingOptionalFields(t *testing.T) {
	payload := []byte(`{"at":"2026-04-20T10:00:00Z","action":"repo.delete","target":"repo1"}`)
	var ev rawEvent
	if err := json.Unmarshal(payload, &ev); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if ev.Action != "repo.delete" {
		t.Errorf("Action = %q, want %q", ev.Action, "repo.delete")
	}
	if ev.ActorUserID != "" {
		t.Errorf("ActorUserID should be empty when omitted, got %q", ev.ActorUserID)
	}
	if ev.ActorIP != "" {
		t.Errorf("ActorIP should be empty when omitted, got %q", ev.ActorIP)
	}
	if ev.Details != nil {
		t.Errorf("Details should be nil when omitted, got %v", ev.Details)
	}
}

func TestRawEvent_JSONDecode_InvalidJSON(t *testing.T) {
	payload := []byte(`{invalid`)
	var ev rawEvent
	if err := json.Unmarshal(payload, &ev); err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestRawEvent_DetailsRoundTrip(t *testing.T) {
	original := rawEvent{
		At:          time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC),
		ActorUserID: "u42",
		ActorIP:     "192.168.1.1",
		Action:      "member.add",
		Target:      "org1",
		Details:     map[string]any{"role": "admin", "teams": []any{"backend", "frontend"}},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var decoded rawEvent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if decoded.Action != original.Action {
		t.Errorf("Action round-trip: got %q, want %q", decoded.Action, original.Action)
	}
	if decoded.ActorUserID != original.ActorUserID {
		t.Errorf("ActorUserID round-trip: got %q, want %q", decoded.ActorUserID, original.ActorUserID)
	}
	if decoded.Target != original.Target {
		t.Errorf("Target round-trip: got %q, want %q", decoded.Target, original.Target)
	}
	if !decoded.At.Equal(original.At) {
		t.Errorf("At round-trip: got %v, want %v", decoded.At, original.At)
	}
}
