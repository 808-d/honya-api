package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/808-d/honya-api/internal/shared/apperror"
	"github.com/808-d/honya-api/pkg/jwt"
	"github.com/808-d/honya-api/pkg/mailer"
	"github.com/808-d/honya-api/pkg/password"
	"github.com/808-d/honya-api/pkg/verification"
	"github.com/jackc/pgx/v5"
)

const (
	purposeVerify = "verify"
	purposeReset  = "reset"
	defaultRole   = "customer"
)

type Service struct {
	repo       *Repository
	tokens     *jwt.Manager
	codes      *verification.Store
	mailer     mailer.Mailer
	appBaseURL string
	webBaseURL string
}

func NewService(
	repo *Repository,
	tokens *jwt.Manager,
	codes *verification.Store,
	mail mailer.Mailer,
	appBaseURL, webBaseURL string,
) *Service {
	return &Service{
		repo:       repo,
		tokens:     tokens,
		codes:      codes,
		mailer:     mail,
		appBaseURL: strings.TrimRight(appBaseURL, "/"),
		webBaseURL: strings.TrimRight(webBaseURL, "/"),
	}
}

// Login verifies credentials and issues an access token. It panics with an
// *apperror.Error on any failure, so callers do not check a returned error.
func (s *Service) Login(ctx context.Context, req LoginRequest) *LoginResult {
	user, err := s.repo.FindByIdentifier(ctx, req.Identifier)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.Unauthorized("invalid credentials"))
		}
		panic(err)
	}

	if !password.Compare(user.HashPassword, req.Password) {
		panic(apperror.Unauthorized("invalid credentials"))
	}

	if !user.IsVerified {
		panic(apperror.Forbidden("email not verified"))
	}

	return s.issueResult(ctx, user)
}

// Register creates an unverified user, stores a verification token in Redis and
// emails a verification link. The account is activated by VerifyRegister.
func (s *Service) Register(ctx context.Context, req RegisterRequest) *RegisterResponse {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		panic(apperror.BadRequest("username, email and password are required"))
	}

	emailExists, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil {
		panic(err)
	}
	if emailExists {
		panic(apperror.Conflict("email already registered"))
	}

	usernameExists, err := s.repo.UsernameExists(ctx, req.Username)
	if err != nil {
		panic(err)
	}
	if usernameExists {
		panic(apperror.Conflict("username already taken"))
	}

	hash, err := password.Hash(req.Password)
	if err != nil {
		panic(err)
	}

	user, err := s.repo.Create(ctx, req.Username, req.DisplayName, hash, req.Email)
	if err != nil {
		panic(err)
	}

	token, err := verification.NewToken()
	if err != nil {
		panic(err)
	}
	if err := s.codes.Save(ctx, purposeVerify, user.PublicID, token); err != nil {
		panic(err)
	}

	link := fmt.Sprintf("%s/api/v1/auth/register/verify?id=%s&token=%s",
		s.appBaseURL, user.PublicID, token)
	if err := s.mailer.Send(ctx, user.Email, "Verify your honya account",
		"Confirm your email by opening this link: "+link); err != nil {
		panic(err)
	}

	return &RegisterResponse{
		ID:       user.PublicID,
		Username: user.Username,
		Email:    user.Email,
	}
}

// VerifyRegister consumes the verification token, activates the account,
// assigns the default role and logs the user in.
func (s *Service) VerifyRegister(ctx context.Context, publicID, token string) *LoginResult {
	if publicID == "" || token == "" {
		panic(apperror.BadRequest("id and token are required"))
	}

	ok, err := s.codes.Take(ctx, purposeVerify, publicID, token)
	if err != nil {
		panic(err)
	}
	if !ok {
		panic(apperror.BadRequest("invalid or expired verification token"))
	}

	user, err := s.repo.Verify(ctx, publicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic(apperror.BadRequest("user not found"))
		}
		panic(err)
	}

	if err := s.repo.AssignRoleByName(ctx, user.ID, defaultRole); err != nil {
		panic(err)
	}

	return s.issueResult(ctx, user)
}

// ForgotPassword emails a reset link to a verified user. It never reveals
// whether the email is registered.
func (s *Service) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) {
	if req.Email == "" {
		panic(apperror.BadRequest("email is required"))
	}

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return
		}
		panic(err)
	}
	if !user.IsVerified {
		return
	}

	token, err := verification.NewToken()
	if err != nil {
		panic(err)
	}
	if err := s.codes.Save(ctx, purposeReset, user.PublicID, token); err != nil {
		panic(err)
	}

	link := fmt.Sprintf("%s/reset-password?id=%s&token=%s",
		s.webBaseURL, user.PublicID, token)
	if err := s.mailer.Send(ctx, user.Email, "Reset your honya password",
		"Reset your password by opening this link: "+link); err != nil {
		panic(err)
	}
}

// ResetPassword consumes the reset token and replaces the password.
func (s *Service) ResetPassword(ctx context.Context, req ResetPasswordRequest) {
	if req.ID == "" || req.Token == "" || req.NewPassword == "" {
		panic(apperror.BadRequest("id, token and new_password are required"))
	}

	ok, err := s.codes.Take(ctx, purposeReset, req.ID, req.Token)
	if err != nil {
		panic(err)
	}
	if !ok {
		panic(apperror.BadRequest("invalid or expired reset token"))
	}

	hash, err := password.Hash(req.NewPassword)
	if err != nil {
		panic(err)
	}
	if err := s.repo.UpdatePassword(ctx, req.ID, hash); err != nil {
		panic(err)
	}
}

// issueResult loads roles and mints an access token for the user.
func (s *Service) issueResult(ctx context.Context, user *User) *LoginResult {
	roles, err := s.repo.RolesByUserID(ctx, user.ID)
	if err != nil {
		panic(err)
	}

	token, expiresAt, err := s.tokens.GenerateAccessToken(user.PublicID, user.Username, roles)
	if err != nil {
		panic(err)
	}

	return &LoginResult{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		User: UserResponse{
			ID:          user.PublicID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Email:       user.Email,
			Phone:       user.Phone,
			Roles:       roles,
		},
	}
}
