package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/helixgitpx/helixgitpx/services/billing-service/internal/provider"
)

func TestPlanValid(t *testing.T) {
	for _, p := range []PlanName{PlanFree, PlanTeam, PlanScale, PlanEnt} {
		if !p.Valid() {
			t.Fatalf("plan %s should be valid", p)
		}
	}
	if PlanName("bogus").Valid() {
		t.Fatal("bogus should be invalid")
	}
}

func TestUpgradePlanRejectsInvalid(t *testing.T) {
	u := &UseCases{Prov: &provider.Stripe{}}
	_, err := u.UpgradePlan(context.Background(), "org1", "sub1", PlanName("super-duper"))
	if err == nil {
		t.Fatal("expected error for invalid plan")
	}
}

func TestUpgradePlanAccepted(t *testing.T) {
	u := &UseCases{Prov: &provider.Stripe{}}
	sub, err := u.UpgradePlan(context.Background(), "org1", "sub1", PlanTeam)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sub.Plan != "team" {
		t.Fatalf("want plan=team got %q", sub.Plan)
	}
}

type spyProvider struct {
	calledMethod string
	calledSubID  string
	calledPlan   string
	err          error
}

func (s *spyProvider) UpsertCustomer(_ context.Context, _, _ string) (provider.Customer, error) {
	return provider.Customer{}, nil
}
func (s *spyProvider) CreateSubscription(_ context.Context, _, _ string) (provider.Subscription, error) {
	return provider.Subscription{}, nil
}
func (s *spyProvider) CancelSubscription(_ context.Context, subID string) error {
	s.calledMethod = "CancelSubscription"
	s.calledSubID = subID
	return s.err
}
func (s *spyProvider) ChangePlan(_ context.Context, subID, plan string) (provider.Subscription, error) {
	s.calledMethod = "ChangePlan"
	s.calledSubID = subID
	s.calledPlan = plan
	return provider.Subscription{ExternalID: subID, Plan: plan, Status: "active"}, s.err
}

func TestUpgradePlan_PassesCorrectPlanToProvider(t *testing.T) {
	spy := &spyProvider{}
	u := &UseCases{Prov: spy}
	_, err := u.UpgradePlan(context.Background(), "org1", "sub_99", PlanScale)
	if err != nil {
		t.Fatal(err)
	}
	if spy.calledMethod != "ChangePlan" {
		t.Errorf("called %q, want ChangePlan", spy.calledMethod)
	}
	if spy.calledSubID != "sub_99" {
		t.Errorf("subID = %q, want %q", spy.calledSubID, "sub_99")
	}
	if spy.calledPlan != "scale" {
		t.Errorf("plan = %q, want %q", spy.calledPlan, "scale")
	}
}

func TestUpgradePlan_PropagatesProviderError(t *testing.T) {
	spy := &spyProvider{err: errors.New("stripe: 402")}
	u := &UseCases{Prov: spy}
	_, err := u.UpgradePlan(context.Background(), "org1", "sub1", PlanTeam)
	if err == nil {
		t.Fatal("expected error from provider")
	}
	if err.Error() != "stripe: 402" {
		t.Errorf("err = %q, want %q", err.Error(), "stripe: 402")
	}
}

func TestCancelPlan_PassesSubID(t *testing.T) {
	spy := &spyProvider{}
	u := &UseCases{Prov: spy}
	if err := u.CancelPlan(context.Background(), "sub_42"); err != nil {
		t.Fatal(err)
	}
	if spy.calledMethod != "CancelSubscription" {
		t.Errorf("called %q, want CancelSubscription", spy.calledMethod)
	}
	if spy.calledSubID != "sub_42" {
		t.Errorf("subID = %q, want %q", spy.calledSubID, "sub_42")
	}
}

func TestCancelPlan_PropagatesProviderError(t *testing.T) {
	spy := &spyProvider{err: errors.New("stripe: 404")}
	u := &UseCases{Prov: spy}
	err := u.CancelPlan(context.Background(), "sub1")
	if err == nil {
		t.Fatal("expected error from provider")
	}
	if err.Error() != "stripe: 404" {
		t.Errorf("err = %q, want %q", err.Error(), "stripe: 404")
	}
}
