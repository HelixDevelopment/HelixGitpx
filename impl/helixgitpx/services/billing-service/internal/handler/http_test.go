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
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("want application/json got %q", ct)
	}
	var body struct {
		Plans []string `json:"plans"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body.Plans) != 4 {
		t.Fatalf("want 4 plans, got %v", body.Plans)
	}
	want := map[string]bool{"free": true, "team": true, "scale": true, "enterprise": true}
	for _, p := range body.Plans {
		if !want[p] {
			t.Fatalf("unexpected plan %q in %v", p, body.Plans)
		}
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
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("want application/json got %q", ct)
	}
	var body struct {
		SubscriptionID string `json:"subscription_id"`
		Plan           string `json:"plan"`
		Status         string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.SubscriptionID == "" {
		t.Fatal("subscription_id should not be empty")
	}
	if body.Plan != "team" {
		t.Fatalf("want plan=team got %q", body.Plan)
	}
	if body.Status == "" {
		t.Fatal("status should not be empty")
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
	var errBody struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&errBody); err != nil {
		t.Fatalf("decode error body: %v", err)
	}
	if errBody.Code != "upgrade_failed" {
		t.Fatalf("want code=upgrade_failed got %q", errBody.Code)
	}
}

func TestUpgradePlan_InvalidJSON(t *testing.T) {
	srv := setup(t)
	resp, err := http.Post(srv.URL+"/v1/subscriptions/sub-1/upgrade", "application/json",
		strings.NewReader(`{bad`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Fatalf("want 400 got %d", resp.StatusCode)
	}
	var errBody struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&errBody); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if errBody.Code != "invalid_json" {
		t.Fatalf("want code=invalid_json got %q", errBody.Code)
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

func TestHealthz_ReturnsStatusOK(t *testing.T) {
	srv := setup(t)
	resp, err := http.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("want application/json content-type, got %q", ct)
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Status != "ok" {
		t.Fatalf("want status=ok got %q", body.Status)
	}
}
