package domain_test

import (
	"testing"

	"github.com/helixgitpx/helixgitpx/services/orgteam/internal/domain"
)

func TestDetectCycle(t *testing.T) {
	// 1 → 2 → 3 (child → parent). Setting 1.parent_id = 3 would cycle.
	parents := map[string]string{
		"2": "1",
		"3": "2",
	}
	if !domain.DetectCycle(parents, "1", "3") {
		t.Errorf("expected cycle for 1 ← 3")
	}
	if domain.DetectCycle(parents, "4", "3") {
		t.Errorf("no cycle for fresh parent")
	}
}

func TestDetectCycle_SelfParent(t *testing.T) {
	if !domain.DetectCycle(map[string]string{}, "1", "1") {
		t.Errorf("expected cycle for self-parent")
	}
}
