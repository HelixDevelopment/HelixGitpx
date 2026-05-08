package adapter

import "testing"

func TestProvider_Constants(t *testing.T) {
	providers := map[Provider]string{
		GitHub: "github",
		GitLab: "gitlab",
		Gitea:  "gitea",
	}
	for p, want := range providers {
		if string(p) != want {
			t.Errorf("Provider %q = %q, want %q", p, p, want)
		}
	}
}

func TestSource_Fields(t *testing.T) {
	s := Source{
		Provider: GitHub,
		BaseURL:  "https://api.github.com",
		Token:    "ghp_test",
		Owner:    "helixgitpx",
		Repo:     "core",
	}
	if s.Provider != GitHub {
		t.Errorf("Provider = %q, want %q", s.Provider, GitHub)
	}
	if s.BaseURL != "https://api.github.com" {
		t.Errorf("BaseURL = %q", s.BaseURL)
	}
	if s.Owner != "helixgitpx" {
		t.Errorf("Owner = %q", s.Owner)
	}
	if s.Repo != "core" {
		t.Errorf("Repo = %q", s.Repo)
	}
}

func TestRefUpdate_Fields(t *testing.T) {
	r := RefUpdate{Name: "refs/heads/main", OldSHA: "abc123", NewSHA: "def456"}
	if r.Name != "refs/heads/main" {
		t.Errorf("Name = %q", r.Name)
	}
	if r.OldSHA != "abc123" {
		t.Errorf("OldSHA = %q", r.OldSHA)
	}
}

func TestPullRequest_Fields(t *testing.T) {
	pr := PullRequest{Number: 42, URL: "https://github.com/org/repo/pull/42"}
	if pr.Number != 42 {
		t.Errorf("Number = %d, want 42", pr.Number)
	}
	if pr.URL == "" {
		t.Error("URL is empty")
	}
}

func TestWebhook_Fields(t *testing.T) {
	wh := Webhook{ID: "wh_1", URL: "http://hook", Events: []string{"push"}}
	if wh.ID != "wh_1" {
		t.Errorf("ID = %q", wh.ID)
	}
	if len(wh.Events) != 1 || wh.Events[0] != "push" {
		t.Errorf("Events = %v", wh.Events)
	}
}
