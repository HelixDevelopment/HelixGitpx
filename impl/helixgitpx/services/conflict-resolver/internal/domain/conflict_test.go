package domain

import (
	"errors"
	"testing"
)

func TestClassify(t *testing.T) {
	cases := []struct {
		name                              string
		refs, labels, rename, meta        bool
		want                              Kind
	}{
		{"refs diverge wins", true, true, true, true, KindRefDivergence},
		{"rename without refs", false, false, true, false, KindRenameCollision},
		{"labels only", false, true, false, false, KindLabelRace},
		{"meta only", false, false, false, true, KindMetaDrift},
		{"nothing", false, false, false, false, KindUnspecified},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Classify(c.refs, c.labels, c.rename, c.meta)
			if got != c.want {
				t.Fatalf("want %v got %v", c.want, got)
			}
		})
	}
}

func TestCanPropose(t *testing.T) {
	if !CanPropose(StatusOpen) || !CanPropose(StatusProposed) {
		t.Fatal("open/proposed should allow proposals")
	}
	if CanPropose(StatusResolved) || CanPropose(StatusRejected) {
		t.Fatal("terminal statuses must forbid proposals")
	}
}

func TestTransition(t *testing.T) {
	if err := Transition(StatusOpen, StatusProposed); err != nil {
		t.Fatal(err)
	}
	if err := Transition(StatusProposed, StatusResolved); err != nil {
		t.Fatal(err)
	}
	if err := Transition(StatusProposed, StatusRejected); err != nil {
		t.Fatal(err)
	}
	if err := Transition(StatusResolved, StatusOpen); !errors.Is(err, ErrResolvedImmutable) {
		t.Fatalf("want immutable, got %v", err)
	}
	if err := Transition(StatusOpen, StatusResolved); !errors.Is(err, ErrResolutionNotProposed) {
		t.Fatalf("want not-proposed, got %v", err)
	}
}

func TestValidateRationale(t *testing.T) {
	if err := ValidateRationale("ai", "   "); !errors.Is(err, ErrEmptyRationale) {
		t.Fatal("blank rationale must be rejected")
	}
	if err := ValidateRationale("human", "LGTM after review"); err != nil {
		t.Fatal(err)
	}
}

func TestKindString(t *testing.T) {
	if KindRefDivergence.String() != "ref_divergence" {
		t.Fatal("string label drift")
	}
}
