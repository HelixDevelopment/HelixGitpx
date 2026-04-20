package domain

import "time"

// Repo mirrors repo.repos.
type Repo struct {
	ID            string
	OrgID         string
	Slug          string
	DefaultBranch string
	LFSEnabled    bool
	CreatedAt     time.Time
}

// Ref mirrors repo.refs.
type Ref struct {
	Name string
	SHA  string
}

// Protection mirrors repo.branch_protections.
type Protection struct {
	RepoID            string
	Pattern           string
	RequireSigned     bool
	RequiredReviewers int
}
