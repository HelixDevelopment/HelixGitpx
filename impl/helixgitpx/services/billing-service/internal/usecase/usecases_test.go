package usecase

import (
	"context"
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
