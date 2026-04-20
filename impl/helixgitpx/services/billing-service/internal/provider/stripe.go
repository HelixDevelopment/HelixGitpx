// Package provider implements billing provider adapters.
package provider

import "context"

type Customer struct {
	ExternalID string
	OrgID      string
	Email      string
}

type Subscription struct {
	ExternalID string
	Plan       string
	Status     string
}

// Provider abstracts the payment processor. Initial impl is Stripe.
type Provider interface {
	UpsertCustomer(ctx context.Context, orgID, email string) (Customer, error)
	CreateSubscription(ctx context.Context, customerID, plan string) (Subscription, error)
	CancelSubscription(ctx context.Context, subscriptionID string) error
	ChangePlan(ctx context.Context, subscriptionID, plan string) (Subscription, error)
}

// Stripe is a placeholder Stripe provider. Real wiring added at GA integration.
type Stripe struct {
	APIKey string
}

func (s *Stripe) UpsertCustomer(ctx context.Context, orgID, email string) (Customer, error) {
	return Customer{ExternalID: "cus_stub_" + orgID, OrgID: orgID, Email: email}, nil
}

func (s *Stripe) CreateSubscription(ctx context.Context, customerID, plan string) (Subscription, error) {
	return Subscription{ExternalID: "sub_stub_" + customerID, Plan: plan, Status: "active"}, nil
}

func (s *Stripe) CancelSubscription(ctx context.Context, subscriptionID string) error {
	return nil
}

func (s *Stripe) ChangePlan(ctx context.Context, subscriptionID, plan string) (Subscription, error) {
	return Subscription{ExternalID: subscriptionID, Plan: plan, Status: "active"}, nil
}
