// Package canonical converts provider-specific webhook payloads into a
// common WebhookEvent shape published to the upstream.webhooks Kafka topic.
package canonical

import (
	"encoding/json"
	"time"
)

// Event is HelixGitpx's provider-agnostic webhook representation. Published
// as JSON to the upstream.webhooks topic; M5 consumers (sync-orchestrator,
// conflict-resolver) decode against this shape.
type Event struct {
	At         time.Time       `json:"at"`
	Provider   string          `json:"provider"`    // "github" | "gitlab" | "gitea"
	DeliveryID string          `json:"delivery_id"`
	EventType  string          `json:"event_type"`
	Repo       string          `json:"repo"`
	Ref        string          `json:"ref,omitempty"`
	BodyRaw    json.RawMessage `json:"body_raw"`
}

// CanonicalizeGitHub maps the GitHub webhook payload + headers to an Event.
func CanonicalizeGitHub(deliveryID, eventType string, body []byte) Event {
	var parsed struct {
		Ref        string `json:"ref"`
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
	}
	_ = json.Unmarshal(body, &parsed)
	return Event{
		At:         time.Now().UTC(),
		Provider:   "github",
		DeliveryID: deliveryID,
		EventType:  eventType,
		Repo:       parsed.Repository.FullName,
		Ref:        parsed.Ref,
		BodyRaw:    body,
	}
}

// CanonicalizeGitLab maps the GitLab webhook payload to an Event.
func CanonicalizeGitLab(deliveryID, eventType string, body []byte) Event {
	var parsed struct {
		Ref     string `json:"ref"`
		Project struct {
			PathWithNamespace string `json:"path_with_namespace"`
		} `json:"project"`
	}
	_ = json.Unmarshal(body, &parsed)
	return Event{
		At:         time.Now().UTC(),
		Provider:   "gitlab",
		DeliveryID: deliveryID,
		EventType:  eventType,
		Repo:       parsed.Project.PathWithNamespace,
		Ref:        parsed.Ref,
		BodyRaw:    body,
	}
}

// CanonicalizeGitea maps the Gitea webhook payload to an Event.
// Same structure as Codeberg and Forgejo (Gitea-compatible).
func CanonicalizeGitea(deliveryID, eventType string, body []byte) Event {
	var parsed struct {
		Ref        string `json:"ref"`
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
	}
	_ = json.Unmarshal(body, &parsed)
	return Event{
		At:         time.Now().UTC(),
		Provider:   "gitea",
		DeliveryID: deliveryID,
		EventType:  eventType,
		Repo:       parsed.Repository.FullName,
		Ref:        parsed.Ref,
		BodyRaw:    body,
	}
}
