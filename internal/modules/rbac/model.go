package rbac

import "time"

// Role is a named permission set. The label is free-form; access is decided by
// the feature codes attached to it.
type Role struct {
	ID          int64
	PublicID    string
	Name        string
	Description string
	Features    []string
	MemberCount int64
	CreateTime  time.Time
	UpdateTime  time.Time
}

// Feature is a grantable capability from the fixed feature catalogue.
type Feature struct {
	ID          int64
	PublicID    string
	Code        string
	Name        string
	Description string
}

// User is a user projection used by role assignment screens.
type User struct {
	ID          int64
	PublicID    string
	Username    string
	DisplayName string
	Email       string
	Phone       *string
	IsVerified  bool
	Roles       []string
	CreateTime  time.Time
}
