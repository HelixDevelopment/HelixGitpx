package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/billing-service/internal/provider"
	"github.com/helixgitpx/helixgitpx/services/billing-service/internal/usecase"
)

func setup(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer((&Handler{
		UseCases: &usecase.UseCases{Prov: &provider.Stripe{}},
	}).Routes())
	t.Cleanup(srv.Close)
	return srv
}

func TestListPlans(t *testing.T) {
	srv := setup(t)
	resp, err := http.Get(srv.URL + "/v1/plans")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	var body struct {
		Plans []string `json:"plans"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if len(body.Plans) != 4 {
		t.Fatalf("want 4 plans, got %v", body.Plans)
	}
}

func TestUpgradePlan(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/subscriptions/sub-1/upgrade", "application/json",
		strings.NewReader(`{"org_id":"o1","plan":"team"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
}

func TestUpgradePlan_BadPlan(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/subscriptions/sub-1/upgrade", "application/json",
		strings.NewReader(`{"org_id":"o1","plan":"nope"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("want 400 got %d", resp.StatusCode)
	}
}

func TestCancel(t *testing.T) {
	srv := setup(t)
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/v1/subscriptions/sub-1/cancel", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		t.Fatalf("want 204 got %d", resp.StatusCode)
	}
}
