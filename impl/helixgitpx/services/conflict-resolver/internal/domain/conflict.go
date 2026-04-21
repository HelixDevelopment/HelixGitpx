// Package domain encodes the conflict-resolver's business rules: when a
// divergence is a conflict, which kind, and what's allowed for resolution.
package domain

import (
	"errors"
	"strings"
)

type Kind int

const (
	KindUnspecified Kind = iota
	KindRefDivergence
	KindLabelRace
	KindRenameCollision
	KindMetaDrift
)

func (k Kind) String() string {
	switch k {
	case KindRefDivergence:
		return "ref_divergence"
	case KindLabelRace:
		return "label_race"
	case KindRenameCollision:
		return "rename_collision"
	case KindMetaDrift:
		return "meta_drift"
	default:
		return "unspecified"
	}
}

type Status int

const (
	StatusUnspecified Status = iota
	StatusOpen
	StatusProposed
	StatusResolved
	StatusRejected
)

var (
	ErrInvalidKind          = errors.New("conflict: invalid kind")
	ErrEmptyRationale       = errors.New("conflict: rationale required")
	ErrResolutionNotProposed = errors.New("conflict: no proposed resolution to accept/reject")
	ErrResolvedImmutable    = errors.New("conflict: resolved conflicts cannot be re-resolved")
)

// Classify maps a (refs-differ, labels-differ, rename, meta) signal tuple
// to the appropriate Kind. The checker always picks the most specific kind;
// Ref divergence wins over label race when both are present.
func Classify(refsDiffer, labelsDiffer, renameCollision, metaDrift bool) Kind {
	switch {
	case refsDiffer:
		return KindRefDivergence
	case renameCollision:
		return KindRenameCollision
	case labelsDiffer:
		return KindLabelRace
	case metaDrift:
		return KindMetaDrift
	default:
		return KindUnspecified
	}
}

// CanPropose gates whether a resolution may be attached to a conflict in
// the given status.
func CanPropose(s Status) bool {
	return s == StatusOpen || s == StatusProposed
}

// Transition validates that `to` is a legal next status for a conflict at
// status `from`.
func Transition(from, to Status) error {
	if from == StatusResolved || from == StatusRejected {
		return ErrResolvedImmutable
	}
	switch to {
	case StatusProposed:
		if from != StatusOpen && from != StatusProposed {
			return ErrResolutionNotProposed
		}
	case StatusResolved, StatusRejected:
		if from != StatusProposed {
			return ErrResolutionNotProposed
		}
	}
	return nil
}

// ValidateRationale enforces that AI and human resolutions explain themselves.
func ValidateRationale(source, rationale string) error {
	_ = source
	if strings.TrimSpace(rationale) == "" {
		return ErrEmptyRationale
	}
	return nil
}
