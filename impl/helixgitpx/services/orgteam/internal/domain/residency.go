package domain

import "errors"

type Residency string

const (
	ResidencyEU Residency = "EU"
	ResidencyUK Residency = "UK"
	ResidencyUS Residency = "US"
)

var ErrInvalidResidency = errors.New("invalid residency zone")

func (r Residency) Valid() bool {
	switch r {
	case ResidencyEU, ResidencyUK, ResidencyUS:
		return true
	}
	return false
}

func SetOrgResidency(currentOwner string, actor string, z Residency) error {
	if currentOwner != actor {
		return errors.New("only owner may change residency")
	}
	if !z.Valid() {
		return ErrInvalidResidency
	}
	return nil
}
