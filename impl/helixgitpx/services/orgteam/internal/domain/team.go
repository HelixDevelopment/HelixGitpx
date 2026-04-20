// Package domain enforces invariants on orgs, nested teams, and memberships.
package domain

// DetectCycle returns true when setting team.parent_id = newParent would
// create a cycle. parents maps child_id → parent_id (existing graph).
// Complexity: O(depth). Returns true if newParent transitively descends
// from team (i.e. newParent is in the subtree rooted at team).
func DetectCycle(parents map[string]string, team, newParent string) bool {
	if team == newParent {
		return true
	}
	cur := newParent
	for i := 0; i < 1_000; i++ {
		if cur == team {
			return true
		}
		next, ok := parents[cur]
		if !ok || next == "" {
			return false
		}
		cur = next
	}
	return true // pathological depth — treat as cycle
}
