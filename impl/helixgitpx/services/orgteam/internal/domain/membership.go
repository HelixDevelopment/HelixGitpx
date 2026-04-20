package domain

// Role matches team.role_enum + the proto Role enum.
type Role string

const (
	RoleViewer Role = "viewer"
	RoleMember Role = "member"
	RoleAdmin  Role = "admin"
	RoleOwner  Role = "owner"
)

// Membership mirrors team.memberships.
type Membership struct {
	ID     string
	TeamID string
	UserID string
	Role   Role
}
