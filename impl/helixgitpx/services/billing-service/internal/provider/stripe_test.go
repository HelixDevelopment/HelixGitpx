package provider

import (
	"context"
	"testing"
)

func TestStripe_UpsertCustomer(t *testing.T) {
	s := &Stripe{APIKey: "test_key"}
	cust, err := s.UpsertCustomer(context.Background(), "org1", "alice@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if cust.ExternalID != "cus_stub_org1" {
		t.Fatalf("want cus_stub_org1 got %q", cust.ExternalID)
	}
	if cust.OrgID != "org1" {
		t.Fatalf("want org1 got %q", cust.OrgID)
	}
	if cust.Email != "alice@example.com" {
		t.Fatalf("want alice@example.com got %q", cust.Email)
	}
}

func TestStripe_CreateSubscription(t *testing.T) {
	s := &Stripe{APIKey: "test_key"}
	sub, err := s.CreateSubscription(context.Background(), "cus_123", "pro")
	if err != nil {
		t.Fatal(err)
	}
	if sub.ExternalID != "sub_stub_cus_123" {
		t.Fatalf("want sub_stub_cus_123 got %q", sub.ExternalID)
	}
	if sub.Plan != "pro" {
		t.Fatalf("want pro got %q", sub.Plan)
	}
	if sub.Status != "active" {
		t.Fatalf("want active got %q", sub.Status)
	}
}

func TestStripe_CancelSubscription(t *testing.T) {
	s := &Stripe{APIKey: "test_key"}
	if err := s.CancelSubscription(context.Background(), "sub_123"); err != nil {
		t.Fatalf("want nil got %v", err)
	}
}

func TestStripe_ChangePlan(t *testing.T) {
	s := &Stripe{APIKey: "test_key"}
	sub, err := s.ChangePlan(context.Background(), "sub_123", "enterprise")
	if err != nil {
		t.Fatal(err)
	}
	if sub.ExternalID != "sub_123" {
		t.Fatalf("want sub_123 got %q", sub.ExternalID)
	}
	if sub.Plan != "enterprise" {
		t.Fatalf("want enterprise got %q", sub.Plan)
	}
	if sub.Status != "active" {
		t.Fatalf("want active got %q", sub.Status)
	}
}

func TestProvider_Interface(t *testing.T) {
	var _ Provider = (*Stripe)(nil)
}
