package provider

import (
	"context"
	"testing"
)

func TestStripe_UpsertCustomer_PropagatesInputs(t *testing.T) {
	s := &Stripe{APIKey: "test_key"}
	cust, err := s.UpsertCustomer(context.Background(), "org1", "alice@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if cust.OrgID != "org1" {
		t.Errorf("OrgID = %q, want %q", cust.OrgID, "org1")
	}
	if cust.Email != "alice@example.com" {
		t.Errorf("Email = %q, want %q", cust.Email, "alice@example.com")
	}
	if cust.ExternalID == "" {
		t.Error("ExternalID is empty")
	}
}

func TestStripe_CreateSubscription_ReturnsActive(t *testing.T) {
	s := &Stripe{APIKey: "test_key"}
	sub, err := s.CreateSubscription(context.Background(), "cus_123", "pro")
	if err != nil {
		t.Fatal(err)
	}
	if sub.Status != "active" {
		t.Errorf("Status = %q, want %q", sub.Status, "active")
	}
	if sub.Plan != "pro" {
		t.Errorf("Plan = %q, want %q", sub.Plan, "pro")
	}
	if sub.ExternalID == "" {
		t.Error("ExternalID is empty")
	}
}

func TestStripe_CancelSubscription_ReturnsNil(t *testing.T) {
	s := &Stripe{APIKey: "test_key"}
	if err := s.CancelSubscription(context.Background(), "sub_123"); err != nil {
		t.Fatalf("CancelSubscription returned error: %v", err)
	}
}

func TestStripe_ChangePlan_ReflectsNewPlan(t *testing.T) {
	s := &Stripe{APIKey: "test_key"}
	sub, err := s.ChangePlan(context.Background(), "sub_123", "enterprise")
	if err != nil {
		t.Fatal(err)
	}
	if sub.Plan != "enterprise" {
		t.Errorf("Plan = %q, want %q", sub.Plan, "enterprise")
	}
	if sub.ExternalID != "sub_123" {
		t.Errorf("ExternalID = %q, want %q", sub.ExternalID, "sub_123")
	}
}

func TestProvider_Interface(t *testing.T) {
	var _ Provider = (*Stripe)(nil)
}

func TestStripe_UpsertCustomer_EmptyInputs(t *testing.T) {
	s := &Stripe{}
	cust, err := s.UpsertCustomer(context.Background(), "", "")
	if err != nil {
		t.Fatal(err)
	}
	if cust.OrgID != "" {
		t.Errorf("OrgID = %q, want empty", cust.OrgID)
	}
}

func TestStripe_CreateSubscription_AllPlans(t *testing.T) {
	s := &Stripe{}
	plans := []string{"free", "team", "scale", "enterprise"}
	for _, plan := range plans {
		sub, err := s.CreateSubscription(context.Background(), "cus_1", plan)
		if err != nil {
			t.Fatalf("plan %s: %v", plan, err)
		}
		if sub.Plan != plan {
			t.Errorf("plan %s: got %q", plan, sub.Plan)
		}
	}
}
