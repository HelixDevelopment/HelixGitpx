// Package usecase wires billing provider actions into audited, policy-checked flows.
package usecase

import (
	"context"
	"errors"

	"github.com/helixgitpx/helixgitpx/services/billing-service/internal/provider"
)

type PlanName string

const (
	PlanFree    PlanName = "free"
	PlanTeam    PlanName = "team"
	PlanScale   PlanName = "scale"
	PlanEnt     PlanName = "enterprise"
)

func (p PlanName) Valid() bool {
	switch p {
	case PlanFree, PlanTeam, PlanScale, PlanEnt:
		return true
	}
	return false
}

type UseCases struct {
	Prov provider.Provider
}

func (u *UseCases) UpgradePlan(ctx context.Context, orgID, subID string, plan PlanName) (provider.Subscription, error) {
	if !plan.Valid() {
		return provider.Subscription{}, errors.New("invalid plan")
	}
	return u.Prov.ChangePlan(ctx, subID, string(plan))
}

func (u *UseCases) CancelPlan(ctx context.Context, subID string) error {
	return u.Prov.CancelSubscription(ctx, subID)
}
