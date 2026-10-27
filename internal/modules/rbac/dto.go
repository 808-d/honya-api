package rbac

// RoleInput is the body of role create and update.
type RoleInput struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	FeatureIDs  []string `json:"feature_ids"`
}

// RoleResponse is the public representation of a role.
type RoleResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Features    []string `json:"features"`
	MemberCount int64    `json:"member_count"`
	CreateTime  string   `json:"create_time"`
	UpdateTime  string   `json:"update_time"`
}

// FeatureResponse is the public representation of a feature.
type FeatureResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UserResponse is a user with their role labels.
type UserResponse struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	DisplayName string   `json:"display_name"`
	Email       string   `json:"email"`
	Phone       *string  `json:"phone,omitempty"`
	IsVerified  bool     `json:"is_verified"`
	Roles       []string `json:"roles"`
	CreateTime  string   `json:"create_time"`
}

// UserListResult is a page of users.
type UserListResult struct {
	Items   []UserResponse `json:"items"`
	Total   int64          `json:"total"`
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
}

// UpdateUserRolesRequest is the body of PUT /admin/users/{id}/roles.
type UpdateUserRolesRequest struct {
	RoleIDs []string `json:"role_ids"`
}
