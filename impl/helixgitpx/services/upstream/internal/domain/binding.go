// Package domain encodes upstream-service invariants: binding direction
// semantics, provider registry, and URL shape validation.
package domain

import (
	"errors"
	"net/url"
	"strings"
)

type Direction int

const (
	DirectionUnspecified Direction = iota
	DirectionReadOnly              // pull from upstream, never push
	DirectionWrite                 // push to upstream; inbound webhooks may still arrive
	DirectionBidirectional         // full federation (read + write + webhooks)
)

func (d Direction) Allows(op Op) bool {
	switch op {
	case OpRead:
		return d == DirectionReadOnly || d == DirectionBidirectional || d == DirectionWrite
	case OpWrite:
		return d == DirectionWrite || d == DirectionBidirectional
	case OpReceiveWebhook:
		return d == DirectionWrite || d == DirectionBidirectional
	default:
		return false
	}
}

type Op int

const (
	OpRead Op = iota
	OpWrite
	OpReceiveWebhook
)

var (
	ErrEmptyRepoID   = errors.New("upstream: repo id is empty")
	ErrInvalidURL    = errors.New("upstream: URL not well-formed")
	ErrUnknownProvider = errors.New("upstream: unknown provider")
	ErrInvalidDirection = errors.New("upstream: invalid direction")
)

// Providers is the fixed set of provider IDs at GA. New providers require
// an ADR and a code change to adapter-pool.
var Providers = map[string]struct{}{
	"github":       {},
	"gitlab":       {},
	"gitea":        {},
	"gitee":        {},
	"bitbucket":    {},
	"gitflic":      {},
	"gitverse":     {},
	"azuredevops":  {},
	"awscodecommit": {},
	"forgejo":      {},
	"sourcehut":    {},
	"generic":      {},
}

// BindingInput is the validated form of an upstream binding request.
type BindingInput struct {
	RepoID    string
	Provider  string
	RawURL    string
	Direction Direction
}

// Validate enforces shape before persistence.
func Validate(in BindingInput) error {
	if strings.TrimSpace(in.RepoID) == "" {
		return ErrEmptyRepoID
	}
	if _, ok := Providers[in.Provider]; !ok {
		return ErrUnknownProvider
	}
	if in.Direction == DirectionUnspecified {
		return ErrInvalidDirection
	}
	u, err := url.Parse(in.RawURL)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "ssh") {
		return ErrInvalidURL
	}
	return nil
}
