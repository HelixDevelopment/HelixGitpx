// Package domain holds repo-service business logic: repos, refs, branch
// protection, and the rules deciding whether a push is allowed.
package domain

import "strings"

// MatchesPattern reports whether a ref short name or full refname matches a
// branch-protection glob pattern. Supports: exact match, '*' wildcard, and
// 'prefix/*' wildcards.
func MatchesPattern(pattern, refName string) bool {
	short := strings.TrimPrefix(refName, "refs/heads/")
	if pattern == "*" {
		return true
	}
	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*") + "/"
		return strings.HasPrefix(short, prefix)
	}
	return pattern == short
}
