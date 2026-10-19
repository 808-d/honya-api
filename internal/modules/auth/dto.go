package auth

// LoginRequest is the body of POST /auth/login.
type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

// UserResponse is the public representation of the authenticated user.
type UserResponse struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	DisplayName string   `json:"display_name"`
	Email       string   `json:"email"`
	Phone       *string  `json:"phone,omitempty"`
	Roles       []string `json:"roles"`
}

// LoginResult is the internal result of a successful login. The access token
// is placed in an httpOnly cookie by the handler and never returned in the body.
type LoginResult struct {
	AccessToken string
	ExpiresAt   int64
	User        UserResponse
}

// LoginResponse is the JSON body returned on a successful login.
type LoginResponse struct {
	ExpiresAt int64        `json:"expires_at"`
	User      UserResponse `json:"user"`
}

// MeResponse is the credential set decoded from an access token.
type MeResponse struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	DisplayName string   `json:"display_name"`
	Email       string   `json:"email"`
	Phone       *string  `json:"phone,omitempty"`
	Roles       []string `json:"roles"`
}

// RegisterRequest is the body of POST /auth/register.
type RegisterRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

// RegisterResponse is the JSON body returned on registration. The account is
// created unverified; a verification link is sent to the email address.
type RegisterResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// ForgotPasswordRequest is the body of POST /auth/password/forgot.
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest is the body of POST /auth/password/reset.
type ResetPasswordRequest struct {
	ID          string `json:"id"`
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}
